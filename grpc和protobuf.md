# grpc和protobuf

## 什么是rpc
RPC（Remote Procedure Call），翻译成中文叫远程过程调用。其设计思路是程序A可以像支持本地函数一样调用远程程序B的一个函数。程序A和程序B很可能不在一台机器上，这中间就需要网络通讯，一般都是基于tcp或者http协议。函数有入参和返回值这里就需要约定一直序列化方式，比较常见的有binary，json，xml。函数的调用过程可以同步和异步。随着分布式程序近几年的大行其道，rpc作为其主要的通讯方式也越来越进入我们的视野，比如GRPC，Thrift。

## go标准库的rpc
golang的标准库就支持rpc。

server端
```golang
package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
)
type Param struct {
	A,B int
}
type Server struct {}

func (t *Server) Multiply(args *Param, reply *int) error {
	*reply = args.A * args.B
	return nil
}
func (t *Server) Divide(args *Param, reply *int) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	*reply = args.A / args.B
	return nil
}

func main() {
	rpc.RegisterName("test", new(Server))
	listener, err := net.Listen("tcp", ":9700")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	for  {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go rpc.ServeConn(conn)
	}
}
```

client端
```golang
package main
import (
	"fmt"
	"log"
	"net/rwida/micro-server
	"time"wida/micro-server
)
type Param struct {
	A,B int
}
func main() {
	client, err := rpc.Dial("tcp", "localhost:9700")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	//同步调用
	var reply int
	err = client.Call("test.Multiply", Param{34,35}, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
	//异步调用
	done := make(chan *rpc.Call, 1)
	client.Go("test.Divide", Param{34,17}, &reply,done)
	select {
		case d := <-done:
			fmt.Println(* d.Reply.(*int))

		case <-time.After(3e9):
			fmt.Println("time out")
	}
}
```

```bash
$ go run server/main.go &
$ go run client/main.go
1190
2
```
程序中我们没有看到参数的序列化和反序列化过程，实际上rpc使用`encoding/gob`做了序列化和反序列化操作，但是`encoding/gob`只支持go语言内部做交互，如果需要夸语言的话就不能用`encoding/gob`了。我们还可以使用标准库中的jsonrpc.

server
```golang
package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Param struct {
	A,B int
}
type Server struct {}

func (t *Server) Multiply(args *Param, reply *int) error {
	*reply = args.A * args.B
	return nil
}
func (t *Server) Divide(args *Param, reply *int) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	*reply = args.A / args.B
	return nil
}

func main() {
	rpc.RegisterName("test", new(Server))
	listener, err := net.Listen("tcp", ":9700")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	for  {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
```

```bash
$ go run server/main.go &
$ echo -e '{"method":"test.Multiply","params":[{"A":34,"B":35}],"id":0}' | nc localhost 9700
{"id":0,"result":1190,"error":null}
$ echo -e '{"method":"test.Divide","params":[{"A":34,"B":17}],"id":1}' | nc localhost 9700
{"id":0,"result":1190,"error":null}
```
这个例子中我们用go写了服务端，客户端我们使用nc直接和服务端交互。可以看到交互的序列化方式是json，因此其它语言也很容易实现和该服务端的交互。

当然我们上面演示的是基于tcp的rpc方式，标准库同时也支持http协议的rpc，感兴趣的同学可以去了解下。

## GRPC
上面我们发了一些篇幅介绍了标准库的rpc，主要的目的是想介绍rpc是什么东西。实际的开发中我们反而比较少用标准库的rpc，我们通常会选择grpc，grpc不仅仅在go生态里头非常流行，在其他语言生态里头同样也非常流行。

grpc是google公司开发的基于http2协议设计和protobuf开发的，高性能，跨语言的rpc框架。这边着重介绍下http2的一个重要特性当tcp多路复用功能，说得直白点就是一个在tcp链接执行多个请求，所以Service A提供N多个服务，Client B和A的所有交互都只用一个链接，这样子可以省很多的链接资源。对tcp连接复用（ connection multiplexing）感兴趣的同学可以阅读下[yamux](https://github.com/hashicorp/yamux)的源码。

### protobuf

protobuf 是google开发的一种平台中立，语言中立，可拓展的数据描述语言。类似json，xml等数据描述语言类似，proto也实现了自己的序列化和反序列化方式，对相对于json来说，protobuf序列化效率更高，体积更小。官方地址[protobuf](https://github.com/protocolbuffers/protobuf)有兴趣的同学可以看看。

#### protobuf安装

我们可以从[protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)里头下载到相应的二进制版本安装protobuf。比如
```
protoc-3.10.0-win32.zip
protoc-3.10.0-linux-x86_64.zip
```
加压文件到文件夹，然后将改文件目录添加到环境变量`PATH`中。
```bash
$ protoc --version
libprotoc 3.10.0
```
ok  protobuf就安装成功了。

#### 安装protobuf go插件

```bash
go get -u -v github.com/golang/protobuf/proto
go get -u -v github.com/golang/protobuf/protoc-gen-go
```

### hello world

新建一个go项目
```bash
├── client
│   └── main.go
├── pb
│   ├── gen.sh
│   └── search.proto
└── server
    └── main.go
``` 

protobuf 消息定义文件为 `pb/search.proto`
```protobuf
syntax = "proto3";

package pb;

service Searcher {
    rpc Search (SearchRequest) returns (SearchReply) {}
}

message SearchRequest {
    bool name = 1;
}

message SearchReply {
    string name = 1;
}
```
这个文件我定义了rpc的服务Searcher，里头定了一个接口Search，同时定义了入参和返回值类型。

我们写一个shell脚本`gen.sh`来生成go程序文件
```
#!/bin/bash
PROTOC=`which protoc`
$PROTOC  --go_out=plugins=grpc:. search.proto
```

cd 到`pb`文件夹执行 `gen.sh` 脚本
```bash
$ ./gen.sh
$ cd ../ && tree
├── client
│   └── main.go
├── pb
│   ├── gen.sh
│   ├── search.pb.go
│   └── search.proto
└── server
    └── main.go
```
可以看到我们多了一个 `search.pb.go`的文件。
我们在server文件下写我们服务main.go
```golang
package main
import (
	"log"
	"net"
	"context"
	"github.com/widaT/gorpc_demo/grpc/pb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchReply, error) {
	return &pb.SearchReply{Name:"hello " + in.GetName()}, nil
}
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSearcherServer(s, &server{})
	s.Serve(lis)
}
```

我们在client文件下写我们的客户端main.go
```golang
package main

import (
	"context"
	"fmt"
	"github.com/widaT/gorpc_demo/grpc/pb"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address     = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSearcherClient(conn)
	s := time.Now()
	r, err := c.Search(context.Background(), &pb.SearchRequest{Name: "world"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r , time.Now().Sub(s))
}
```

```bash
$ go run server/main.go &
$ go run client/main.go &
name:"hello world"  14.767013ms
```

### grpc stream
上面的hello world程序演示了grpc的一般用法，这种方式能满足大部分的场景。grpc还提供了双向或者单向流的方式我们成为grpc stream，stream的方式一般用在有大量的数据交互或者长时间交互。

我们修改的grpc服务定义文件`pb/search.proto`,在service增加 `rpc Search2 (stream SearchRequest) returns (stream SearchReply) {}`
```protobuf
syntax = "proto3";
package pb;
service Searcher {
    rpc Search (SearchRequest) returns (SearchReply) {}
    rpc Search2 (stream SearchRequest) returns (stream SearchReply) {}
}
message SearchRequest {
    string name = 1;
}
message SearchReply {
    string name = 1;
}
```

服务端我们修改代码，添加一个`func (s *server) Search2(pb.Searcher_Search2Server) error`的一个实现方法。
```golang
package main

import (
	"io"
	"log"
	"net"

	"context"
	"github.com/widaT/gorpc_demo/grpc/pb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchReply, error) {
	return &pb.SearchReply{Name:"hello " + in.GetName()}, nil
}
func (s *server) Search2(stream pb.Searcher_Search2Server) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		reply := &pb.SearchReply{ Name:"hello:" + args.GetName()}
		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSearcherServer(s, &server{})
	s.Serve(lis)
}
```

为了方便阅读代码，client代码我们重写
```golang
package main

import (
	"context"
	"fmt"
	"github.com/widaT/gorpc_demo/grpc/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSearcherClient(conn)

	stream, err := c.Search2(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	go func() { //这边启动一个goroutine 发送请求
		for {
			if err := stream.Send(&pb.SearchRequest{Name: "world"}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}()

	for { //主goroutine 一直接收结果
		rep, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println(rep.GetName())
	}
}
```

运行代码
```bash
$ go run server/main.go &
$ go run client/main.go 
hello:world
hello:world
hello:world
...
```

# 总结
本文只是对golang使用grpc做了简要的介绍，grpc的使用范围很广，需要我们持续了解和学习。
下面推荐个资料
- [grpc go官方examples](https://github.com/grpc/grpc-go/tree/master/examples)
- [《Go语言高级编程》RPC和Protobuf](https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch4-rpc/readme.html)

# 参考文档

- [grpc go](https://grpc.io/docs/quickstart/go/)