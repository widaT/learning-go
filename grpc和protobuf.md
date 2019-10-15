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
	"net/rpc"
	"time"
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
# go run server/main.go &
# go run client/main.go
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

client 
```golang
package main
import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)
type Param struct {
	A,B int
}
func main() {
	conn, err := net.Dial("tcp", "localhost:9700")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	client :=rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
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