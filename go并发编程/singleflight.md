# singleflight

## 什么是singleflight

`singleflight`是go拓展同步库中实现的一种针对重复的函数调用抑制机制，也就是说一组相同的操作（函数），只有一个函数能被执行，执行后的结果返回给其他同组的函数。`singleflight`应该算是一种并发模型，非常适合当redis某个key缓存失效时，这个时候一堆请求去数据库来取数据然后更新redis缓存，如果我们使用`singlefilght`并发模型的话，那就是redis key失效的时候，一堆去数据库的请求只有一个能成功，它将更新该redis key的value值，同时把value值给其他相同请求。

`singlefilght`中的`Group`结构体

```go
type Group struct {
	mu sync.Mutex       // 互斥锁保护 m
	m  map[string]*call // 存形同key的 函数
}
```

`singleflight`有三个对外api

```go
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)
func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result
func (g *Group) Forget(key string)
```

`Group.Do`和`Group.DoChan`功能类似，`Do`是同步返回，`DoChan`返回一个chan实现异步返回。
`Group.Forget`方法可以通知 `singleflight` 在map删除某个key，接下来对该key的调用就会直接执行方法而不是等待前面的函数返回。

## 使用singleflight

```go
package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"golang.org/x/sync/singleflight"
	"sync"
	"sync/atomic"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379"})

	var g singleflight.Group
	var callTimes int32 =0
	fakeGetData := func() int32 {//模拟去数据库取数据，记录被调用次数
		callTimes = atomic.AddInt32(&callTimes,1)
		return callTimes
	}

	wg := sync.WaitGroup{}
	for i:=0;i<10;i++ { //模拟10并发
		wg.Add(1)
		go func() {
			defer wg.Done()
			ret,err,_:= g.Do("wida", func() (i interface{}, e error) {
				num := fakeGetData()
				client.Set("wida",num,0)
				return num,nil
			})
			fmt.Println(ret,err)
		}()
	}
	wg.Wait()
	fmt.Printf("callTimes %d \n",callTimes)
	ret,_:=client.Get("wida").Int()
	fmt.Printf("redis value %d \n",ret)
}
```

运行结果

```bash
$ go run main.go
1 <nil>
1 <nil>
1 <nil>
1 <nil>
1 <nil>
1 <nil>
1 <nil>
1 <nil>
1 <nil>
1 <nil>
callTimes 1 
redis value 1 
```

我们看到10个并发请求，`fakeGetData`只被调用了一次，reids的值也被设置为1，10个请求拿到了相同的结果。

如果要使用`DoChan`的方式只需要稍微修改下
```go
go func() {
			defer wg.Done()
			retChan:= g.DoChan("wida", func() (i interface{}, e error) {
				num := fakeGetData()
				client.Set("wida",num,0)
				return num,nil
			})
			ret := <- retChan
			fmt.Println(ret)
}()
```        


## 总结

本小节介绍了go拓展同步库中`singleflight`（golang.org/x/sync/singleflight），以及介绍了`singleflight`使用方式和适合的场景。