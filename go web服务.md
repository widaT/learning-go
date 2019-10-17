# go web服务
在绪论中我们简单使用go写了个一个web服务，本文将展开介绍go如何写web服务。我们再看看之前的代码，
```
package main
import (
    "io"
    "net/http"
)
func helloHandler(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "hello, world!\n")
}

func main() {
    http.HandleFunc("/", helloHandler)
    http.ListenAndServe(":8888", nil)
}
```
看起来代码没几行，结构也很清晰， `helloHandler`的函数，它有两个参数`http.ResponseWriter`和`*http.Request`，`http.ResponseWriter` 实际上是一个interface
```golang
type ResponseWriter interface {
	Header() Header        //返回Header map （type Header map[string][]string），用于设置和获取相关的http header信息
	Write([]byte) (int, error)  //返回给client body的内容
	WriteHeader(statusCode int) //返回给client 的http code码
}
```
`http.Request` 包含这次http请求的request信息 入http method，url，body等。

`http.HandleFunc("/", helloHandler)`的作用是使用go定义的默认http路由，添加了 `/` 的请求路径到  `helloHandler`的路由，默认的http路由`DefaultServeMux`，感兴趣的同学可以阅读下源码，这边不展开说。
`http.ListenAndServe` 有两个参数第一个参数是本地地址（ip+port），第二个参数为`nil`就是使用默认的路由器 `DefaultServeMux`。他的作用是创建tcp服务，监听（ip+port），并且创建一个`goroutine` 请处理每一个http请求。


## 自己实现http路由
```golang
package main

import (
	"net/http"
)

type MyServer struct {
	router map[string]http.HandlerFunc
}

func (s *MyServer)ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if  fn,found := s.router[r.URL.Path] ;found {
		fn(rw,r)
		return
	}
	rw.WriteHeader(404)
	rw.Write([]byte("page not found"))
}
func (s *MyServer)Add(path string,fn http.HandlerFunc)  {
	s.router[path] = fn
}
func NewServer() *MyServer  {
	return &MyServer{
		router: make(map[string]http.HandlerFunc),
	}
}
func main()  {
	s := NewServer()
	s.Add("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("hello world"))
	})
	http.ListenAndServe(":8888", s)
}

```
运行一下
```bash
$ go run main.go
$ curl http://localhost:8888
hello world
```

我们写了一堆关于结构体 `MyServer` 的代码，真正把 `MyServer` 和http服务绑定的只有 `http.ListenAndServe(":8888", s) `这一行代码。
说明`http.ListenAndServe`才是我们了解go http服务的关键。我们看下`http.ListenAndServe`的函数定义
`func ListenAndServe(addr string, handler Handler)`，第一个参数我们上文介绍过，第二个参数是一个interface
```golang
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
它定义了`ServeHTTP(ResponseWriter, *Request)`的方法。我们的`MyServer`刚好实现了这个方法
```golang
func (s *MyServer)ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if  fn,found := s.router[r.URL.Path] ;found {  //从 map查找对应的 paht=》http.HandlerFunc 映射
		fn(rw,r)  //真正处理请求，返回消息给client的地方
		return
	}
	rw.WriteHeader(404) //给client发 404
	rw.Write([]byte("page not found"))
}
```
到此我们知道 `MyServer`的 `ServeHTTP` 方法是连接go http server底层和我们写的代码的一个桥梁。我们所需要的每次http请求信息都在`http.Request`中，然后可以通过 `http.ResponseWriter`给客户端回写消息。我们甚至不用关心底层做了什么，我们只需要专注于处理client的每个http请求。