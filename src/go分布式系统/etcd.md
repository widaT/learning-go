# Etcd

在本章节我们不止一次提到etcd，因为etcd这个项目在go语言分布式生态里头非常的重要。etcd个高可用的分布式键值(key-value)数据库，采用raft协议作为一致性算法，基于Go语言实现。etcd在go分布式系统生态里头的位置相当于zookeeper在java分布式系统生态的位置。tidb的raft系统是基于etcd的实现然后用rust实现的，tidb的pd系统还直接内嵌了etcd。在微服务生态里头etcd可以做服务发现组件使用，kubernetes也用etcd保存集群所有的网络配置和对象的状态信息。


## etcd的特点

- 简单：提供定义明确的面向用户的API（有http和grpc实现，注意http和grpc分别用了两套存储，数据不同）
- 安全：支持client端证书的TLS
- 快速：基准测试支持10,000次write/s
- 可靠：使用raft协议保证数据一致性

对于想研究分布式系统的gopher来说，etcd是一个非常值得好好阅读源码的项目。它的raft实现，mvcc，wal实现都非常典型。

## etcd使用

etcd项目有两个go client，一个基于http restful对于源码目录里头的client，另外一个基于grpc对于源码目录的clientv3，需要注意是这两个api操作的数据是不一样的。client支持的功能较小，而且性能也比较有问题，在go生态里头都会用基于grpc的clientv3。


## 连接etcd

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/etcd-io/etcd/clientv3"
)

var client *clientv3.Client

func init() {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"}, //集群有多个需要写多台
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
}
```


### kv操作

etcd 支持来时nosql数据库一样的 kv操作。

```go
//kv操作
func kv() {
	_, err := client.Put(context.Background(), "sample_key", "sample_value")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := clientv3.NewKV(client).Get(context.Background(), "sample_key")
	if n := len(resp.Kvs); n == 0 {
		fmt.Println("not found")
		return
	}
	fmt.Println(string(resp.Kvs[0].Value))

	client.Delete(context.Background(), "sample_key")

	resp, err = clientv3.NewKV(client).Get(context.Background(), "sample_key")
	if n := len(resp.Kvs); n == 0 {
		fmt.Println("not found")
		return
	}
	fmt.Println(string(resp.Kvs[0].Value))
}
```

运行一下

```bash
$ go run main.go
sample_value
not foun
```

### watch（监听）

监听是指etcd可以对单个key或者摸个前缀的key进行监听，如果这些key的信息有变化则会通知监听的client。利用这个功能etcd可以做配置中心，或者做服务发现。

```go
//监听 watch
func watch() {
	go func() {
		for {
			ch := client.Watch(context.Background(), "sample_key") //监听sample_key
			for wresp := range ch {
				for _, ev := range wresp.Events { //打印事件
					fmt.Println(ev) 
				}
			}
		}
	}()
	_, err := client.Put(context.Background(), "sample_key", "sample_value1")
	if err != nil {
		log.Fatal(err)
	}
	client.Delete(context.Background(), "sample_key")
}
```

运行一下

```bash
$ go run main.go
&{PUT key:"sample_key" create_revision:18157 mod_revision:18157 version:1 value:"sample_value1"  <nil> {} [] 0} #PUT事件
&{DELETE key:"sample_key" mod_revision:18158  <nil> {} [] 0}#DELETE事件
```

我看到通过Watch 我们监控到了PUT和DELETE事件，在实际项目中可以根据业务需求针对相关的事件做成处理。


### Lease（租约)

我们在上个小节介绍了 全局时间戳的实现，里头用到了etcd Lease来做选主服务，etcd的Lease可以用来做服务状态监控，再配合watch就能做到故障恢复。上个小节的选主服务本质上是我们常见的主备切换过程。接下来我们看下etcd Lease是怎么使用的。

```go
//租约
func lease() {
	lease := clientv3.NewLease(client)
	//设置10秒租约 （过期时间为10秒）
	leaseResp, err := lease.Grant(context.TODO(), 10)
	if err != nil {
		log.Fatal(err)
	}

	leaseID := leaseResp.ID
	fmt.Printf("ttl %d \n", leaseResp.TTL)

	keepRespChan, err := lease.KeepAlive(context.TODO(), leaseID) //启动自动续租服务
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() { //新开一个协程监控 keepalive结果
	outer:
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约失效了")
					break outer
				} else {

					fmt.Printf("续租成功 time:%d id:%d,ttl:%d \n", time.Now().Unix(), keepResp.ID, keepResp.TTL)
				}

			}
		}
	}()

	kv := clientv3.NewKV(client)
	putResp, err := kv.Put(context.TODO(), "sample_key_lease", "", clientv3.WithLease(leaseID))
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("写入成功", putResp.Header.Revision)
	}
}
```

运行一下

```bash
$ go run main.go
ttl 10 
续租成功 time:1577158783 id:7587839846421449960,ttl:10 
写入成功 18159
续租成功 time:1577158786 id:7587839846421449960,ttl:10 
续租成功 time:1577158790 id:7587839846421449960,ttl:10 
续租成功 time:1577158793 id:7587839846421449960,ttl:10 
续租成功 time:1577158797 id:7587839846421449960,ttl:10 
续租成功 time:1577158800 id:7587839846421449960,ttl:10 
```

我们看到超时时间为10秒，续租程序会在租约失效前大概（1/3 TTL）
```go
	nextKeepAlive := time.Now().Add((time.Duration(karesp.TTL) * time.Second) / 3.0)
```
自动续租（修改TTL），从而实现key一直有效。在服务宕机或者异常的时候没有执行续租，那么这个key将会在10秒后失效，如果其他client watch这个key的话就能监听的这个key被删除了。

## 总结

本小节简要介绍了etcd的一些特性，以及go语言etcd v3版本的api是如何使用的。最后还是推荐想了解分布式系统的gopher去阅读etcd源码。