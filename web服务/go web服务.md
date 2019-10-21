# go web服务

在绪论中我们简单使用go写了个一个web服务，本文将展开介绍go如何写web服务。我们再看看之前的代码。
```golang
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
看起来代码没几行，结构也很清晰， `helloHandler`的函数，它有两个参数`http.ResponseWriter`和`*http.Request`。

- `http.ResponseWriter` 实际上是一个interface
```golang
type ResponseWriter interface {
	Header() Header        //返回Header map （type Header map[string][]string），用于设置和获取相关的http header信息
	Write([]byte) (int, error)  //返回给client body的内容
	WriteHeader(statusCode int) //返回给client 的http code码
}
```
- `http.Request` 包含这次http请求的request信息 入http method，url，body等。

`http.HandleFunc("/", helloHandler)`的作用是注册了 path为`/`处理函数为 `helloHandler`的路由（这边使用了默认的http路由`DefaultServeMux`，感兴趣的同学可以阅读下源码）。
`http.ListenAndServe` 有两个参数第一个参数是本地地址（ip+port），第二个参数为`nil`就是使用默认的路由器 `DefaultServeMux`。它的作用是创建tcp服务，监听（ip+port），并且创建一个`goroutine` 请处理每一个http请求。

接下来我们通过自己实现http路由来加深下go web服务运行方式的了解。
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
说明`http.ListenAndServe`才是我们了解go web服务的关键。

我们看下`http.ListenAndServe`的函数定义
`func ListenAndServe(addr string, handler Handler)`

- 第一个参数我们上文介绍过，
- 第二个参数是一个interface

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

到此我们知道`MyServer`的 `ServeHTTP`方法是连接go http server底层和我们写的代码的一个桥梁。我们所需要的每次http请求信息都在`http.Request`中，然后可以通过 `http.ResponseWriter`给客户端回写消息。

我们甚至不用关心底层做了什么，我们只需要专注于处理client的每个http请求。

接下来我们再改进下代码，定义我们自己的`HandlerFunc`和`Context`，`HandlerFunc`能让我们少写点代码，`Context`的封装则可以定制写我们框架专属的特性，例如`SayHello`方法，代码如下：

```golang
	package main
	import (
		"net/http"
	)
	type HandlerFunc func(*Context)
	type MyServer struct {
		router map[string]map[string]HandlerFunc
	}

	type Context struct {
		Rw  http.ResponseWriter
		R *http.Request
	}
	func (ctx *Context)SayHello()  {
		ctx.Rw.WriteHeader(200)
		ctx.Rw.Write([]byte("hello world"))
	}
	func (s *MyServer)ServeHTTP(rw http.ResponseWriter, r *http.Request) {
		if  _,found := s.router[r.Method] ;found {
			if fn,found :=s.router[r.Method][r.URL.Path];found {
				fn(&Context{Rw:rw,R:r})
				return
			}
		}
		rw.WriteHeader(404)
		rw.Write([]byte("page not found"))
	}
	func (s *MyServer)Get(path string,fn HandlerFunc)  {
		if s.router["GET"] == nil {
			s.router["GET"] = make(map[string]HandlerFunc)
		}
		s.router["GET"][path] = fn
	}
	func (s *MyServer)Post(path string,fn HandlerFunc)  {
		if s.router["POST"] == nil {
			s.router["POST"] = make(map[string]HandlerFunc)
		}
		s.router["POST"][path] = fn
	}
	func NewServer() *MyServer  {
		return &MyServer{
			router: make(map[string]map[string]HandlerFunc),
		}
	}
	func main()  {
		s := NewServer()
		s.Get("/get", func(ctx *Context) {
			ctx.SayHello()
		})

		s.Post("/post", func(ctx *Context) {
			ctx.SayHello()
		})

		http.ListenAndServe(":8888", s)
	}
```

运行一下

```bash
$ go run main.go
$ curl http://localhost:8888/get
hello world
$ curl -d "" http://localhost:8888/post
hello world
```

如果有同学了解过go web框架的话，相信已经感觉到上面的代码和常见框架运用代码已经很像了。
是的，其实go比较流行的web框架比如gin,beego都是利用类似的原理写的，只是他们对`Context`的封装更加丰富，路由使用了树状结构，还有更多的`middleware`。

那么接下来我们看下`httprouter` 和 `middleware`

## http路由（httprouter）
目前大多数流行框架的路由都采用压缩前缀树（compact prefix tree 或者Radix tree），通常每个method都是一颗前缀树。这个部分我们不展开讲，很多主流框架采用方式都比较类似，可以参考[httprouter](https://github.com/julienschmidt/httprouter)项目，做下深入研究。其实对于一些简单场景（例如没有 path 参数，接口很少的场景）根本不需要用树状结构，直接用map可以实现，而且效率更高。


## 中间件（middleware）
先看一段代码
```golang
package main

import (
	"log"
	"net/http"
	"time"
)
type HandlerFunc func(*Context)
type MyServer struct {
	router map[string]map[string]HandlerFunc
}

type Context struct {
	Rw  http.ResponseWriter
	R *http.Request
}

func timeMiddleware(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(ctx *Context) {
		start := time.Now()
		next(ctx)
		elapsed := time.Since(start)
		log.Println("time elapsed",elapsed)
	})
}

func (ctx *Context)SayHello()  {
	ctx.Rw.WriteHeader(200)
	ctx.Rw.Write([]byte("hello world"))
}

func (s *MyServer)ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if  _,found := s.router[r.Method] ;found {
		if fn,found :=s.router[r.Method][r.URL.Path];found {
			fn(&Context{Rw:rw,R:r})
			return
		}
	}
	rw.WriteHeader(404)
	rw.Write([]byte("page not found"))
}
func (s *MyServer)Get(path string,fn HandlerFunc)  {
	if s.router["GET"] == nil {
		s.router["GET"] = make(map[string]HandlerFunc)
	}
	s.router["GET"][path] = fn
}

func NewServer() *MyServer  {
	return &MyServer{
		router: make(map[string]map[string]HandlerFunc),
	}
}
func main()  {
	s := NewServer()
	s.Get("/get", timeMiddleware(func(ctx *Context) {
		ctx.SayHello()
	}))

	http.ListenAndServe(":8888", s)
}
```

```bash
$ go run main.go
2019/10/18 10:46:51 time elapsed 2.257µs  #调用后会出现
$ curl -curl  http://localhost:8888/get
hello world
```

上面的代码 我们定义了`timeMiddleware`函数， 
```golang
func timeMiddleware(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(ctx *Context) {
		start := time.Now()
		next(ctx)
		elapsed := time.Since(start)
		log.Println("time elapsed",elapsed)
	})
}
```
用来包裹一个我们的 `/get` handler 函数。
```golang
	s.Get("/get", timeMiddleware(func(ctx *Context) {
			ctx.SayHello()
		}))
```	

`timeMiddleware`函数的作用是计算`handlerFunc`的耗时， 类似`timeMiddleware`这种函数我们称作中间件。
从上面的代码可以看出，中间件可以不止一层，我们稍微改下代码，使用两层timeMiddleware。
```golang
s.Get("/get", timeMiddleware(timeMiddleware(func(ctx *Context) {
			ctx.SayHello()
		})))
```
运行一下
```bash
$ go run main.go
2019/10/18 11:14:52 time elapsed 2.037µs    #调用后会出现
2019/10/18 11:14:52 time elapsed 70.15µs  
$ curl -curl  http://localhost:8888/get
hello world
```

到这边我们大概了解了中间件的工作原理。我们再对代码做下封装，我们定义了`type MiddleWare func(HandlerFunc)HandlerFunc`，同时使用`MyServer.Use`函数添加中间件，修改了`MyServer.Get`
```golang
package main

import (
	"log"
	"net/http"
	"time"
)
type HandlerFunc func(*Context)
type MiddleWare func(HandlerFunc)HandlerFunc
type MyServer struct {
	router map[string]map[string]HandlerFunc
	chain []MiddleWare
}

type Context struct {
	Rw  http.ResponseWriter
	R *http.Request
}

func timeMiddleware(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(ctx *Context) {
		start := time.Now()
		next(ctx)
		elapsed := time.Since(start)
		log.Println("time elapsed",elapsed)
	})
}

func (ctx *Context)SayHello()  {
	ctx.Rw.WriteHeader(200)
	ctx.Rw.Write([]byte("hello world"))
}

func (s *MyServer)ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if  _,found := s.router[r.Method] ;found {
		if fn,found :=s.router[r.Method][r.URL.Path];found {
			fn(&Context{Rw:rw,R:r})
			return
		}
	}
	rw.WriteHeader(404)
	rw.Write([]byte("page not found"))
}

func (s *MyServer)Use(middleware ...MiddleWare)  {
	for _,m:= range middleware {
		s.chain = append(s.chain,m)
	}
}

func (s *MyServer)Get(path string,fn HandlerFunc)  {
	if s.router["GET"] == nil {
		s.router["GET"] = make(map[string]HandlerFunc)
	}
	handler := fn
	for i := len(s.chain) - 1; i >= 0; i-- {
		handler = s.chain[i](handler)
	}
	s.router["GET"][path] = handler
}

func NewServer() *MyServer  {
	return &MyServer{
		router: make(map[string]map[string]HandlerFunc),
	}
}
func main()  {
	s := NewServer()
	s.Use(timeMiddleware,timeMiddleware)
	s.Get("/get", func(ctx *Context) {
		ctx.SayHello()
	})
	http.ListenAndServe(":8888", s)
}
```

再次运行一下
```bash
$ go run main.go
2019/10/18 11:31:32 time elapsed 1.162µs #调用后会出现
2019/10/18 11:31:32 time elapsed 11.941µs
$ curl -curl  http://localhost:8888/get
hello world
```
ok，代码运行正常。
中间件的思路非常适合做压缩，用户鉴权，access日志，流量控制，安全校验等功能。
注意本文介绍的中间件实现方式和gin的实现方式有点差别，但是核心思路是一样的，个人觉得`gin`的实现方式有点繁琐，有兴趣的同学可以去研究下`gin`中间件实现方式。

# 总结
本文简要的介绍了go web服务的相关编写方式，同时简要介绍了 `httprouter` 和 `middleware` 这两个go web框架核心的组件，希望这些篇幅对大家后续web框架有更深入的理解。

# 参考资料
- [gin框架](https://github.com/gin-gonic/gin)
- [gin的middleware](https://github.com/gin-gonic/contrib)