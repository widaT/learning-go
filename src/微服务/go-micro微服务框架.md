# go-micro微服务框架

## go-micro简介
Go Micro是可插拔的微服务开发框架。micro的设计哲学是可插拔的架构理念，她提供可快速构建系统的组件，并且可以根据自身的需求剥离默认实现并自行定制。详细的介绍可参考官方中文文档[Go Micro](https://micro.mu/),目前go micro这个项目最新版本为v3。

## 安装go-micro依赖
    
-   安装protobuf 参考[grcp和protobuf](./grpc和protobuf.md) 中protobuf的安装方式
-   安装protoc-gen-micro，这个是go-micro定制的protobuf插件，用于生成go-micro定制的类似grpc的代码。
    ```
    go get github.com/micro/protoc-gen-micro
    ```
## 写一个微服务版本的hello world程序

这个程序包含一个service，web，cli：
- service作为微服务的具体体现，他的功能是提供某一类服务，我们这边的稍作简化，只提供一个grpc服务。在实际场景中service是某一类服务接口的集合，接口数量不宜太多，得考虑负载情况和业务情况相结合做适当的拆分整合。
- web作为前端应用，对外提供http服务。在实践场景中web通常会整合各种后端service包装成业务接口提供给终端或者第三方合作伙伴使用。
- cli也是前端应用的另一种形式，不对外提供http服务，可以在命令行下的脚本任务或者某些守护进程。

### 跟grpc一样我们先编写proto文件，定义service

```protobuf
syntax = "proto3";

option go_package = "../proto";

service Say {
	rpc Hello(Request) returns (SayResponse) {}
}

message Request {
	string name = 1;
}

message Pair {
    int32 key = 1;
    string values = 2;
}

message SayResponse {
    string msg = 1;
    // 数组
    repeated string values = 2;
    // map
    map<string, Pair> header = 3;
    RespType type = 4;
}

enum RespType {
    NONE = 0;
    ASCEND = 1;
    DESCEND = 2;
}

```

我们写一个shell脚本来生成golang源码文件

```bash
#!/bin/bash
protoc --go_out=. --micro_out=. test.proto
```

目录结构
```bash
$ tree
.
├── cli.Dockerfile
├── cmd
│   ├── cli
│   │   └── main.go
│   ├── service
│   │   └── main.go
│   └── web
│       └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
├── proto
│   ├── gen.sh
│   ├── test.pb.go
│   ├── test.pb.micro.go
│   └── test.proto
├── README.md
├── service.Dockerfile
└── web.Dockerfile
```
### 编写service
注意go-micro这边使用了默认的mdns来做服务发现，mdns的相关原理可参考[Multicast_DNS](https://en.wikipedia.org/wiki/Multicast_DNS)

```golang
package main

import (
	"context"
	"log"
	hello "mircotest/proto"
	"os"

	"github.com/asim/go-micro/v3"
)

type Hello struct{}

func (s *Hello) Hello(ctx context.Context, req *hello.Request, rsp *hello.SayResponse) error {
	log.Print("Received Say.Hello request")
	hostname, _ := os.Hostname()
	rsp.Msg = "Hello " + req.Name + " ,Im " + hostname
	rsp.Header = make(map[string]*hello.Pair)
	rsp.Header["name"] = &hello.Pair{Key: 1, Values: "abc"}
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("wida.micro.srv.greeter"),
	)
	service.Init()

	// Register Handlers
	hello.RegisterSayHandler(service.Server(), new(Hello))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

```

### 编写web
```golang
package main

import (
	"context"
	"fmt"
	"log"
	hello "mircotest/proto"
	"net/http"

	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/web"
)

func main() {
	service := web.NewService(
		web.Name("wida.micro.web.greeter"),
		web.Address(":8009"),
	)

	service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()
			name := r.Form.Get("name")
			if len(name) == 0 {
				name = "World"
			}
			cl := hello.NewSayService("wida.micro.srv.greeter", client.DefaultClient)
			rsp, err := cl.Hello(context.Background(), &hello.Request{
				Name: name,
			})
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Write([]byte(`<html><body><h1>` + rsp.Msg + `</h1></body></html>`))
			return
		}
		fmt.Fprint(w, `<html><body><h1>Enter Name<h1><form method=post><input name=name type=text /></form></body></html>`)
	})

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

### 编写cli
```golang
package main

import (
	"context"
	"fmt"
	hello "mircotest/proto"

	roundrobin "github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v3"
	"github.com/asim/go-micro/v3"
)

func main() {
	wrapper := roundrobin.NewClientWrapper()
	service := micro.NewService(
		micro.WrapClient(wrapper),
	)
	service.Init()
	cl := hello.NewSayService("wida.micro.srv.greeter", service.Client())
	rsp, err := cl.Hello(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", rsp.Msg)
}

```

### 编写Makefile简化编译过程

```bash
export GOCMD=GO111MODULE=on CGO_ENABLED=0 go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

all : service web cli

service :
	@echo "build service"
	@mkdir -p build
	$(GOBUILD) -a -installsuffix cgo -ldflags '-w' -o build/service cmd/service/*.go

web :
	@echo "build web"
	$(GOBUILD) -a -installsuffix cgo -ldflags '-w' -o build/web cmd/web/*.go

cli :
	@echo "build cli"
	$(GOBUILD) -a -installsuffix cgo -ldflags '-w' -o build/cli cmd/cli/*.go
	
docker:
	@echo "build docker images"
	docker-compose up --build
	
.PHONY: clean
clean:
	@rm -rf build/

.PHONY: proto
proto:
	protoc --go_out=. --micro_out=. test.proto
```    

这边需要注意我们编译golang的时候使用了`CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w' -o` 这样的编译参数是为了去掉编译后的golang可执行文件对cgo的依赖，我们的程序要放在`alpine:latest`容器中，如果依赖cgo则运行不起来。

```
$ make
build service
GO111MODULE=on CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w' -o build/service cmd/service/*.go
build web
GO111MODULE=on CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w' -o build/web cmd/web/*.go
build cli
GO111MODULE=on CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w' -o build/cli cmd/cli/*.go
$ tree 
.
├── build
│   ├── cli
│   ├── service
│   └── web
├── cli.Dockerfile
├── cmd
│   ├── cli
│   │   └── main.go
│   ├── service
│   │   └── main.go
│   └── web
│       └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
├── proto
│   ├── gen.sh
│   ├── test.pb.go
│   ├── test.pb.micro.go
│   └── test.proto
├── README.md
├── service.Dockerfile
└── web.Dockerfile
```

## 微服务运行环境

一般情况下微服务是运行在容器中的，容器化部署会带来很大的运维管理便利性。但是容器化运行不是必选项，通常我们在开发调试的时候我们选择可以物理机运行。等程序运行稳定后，我们才打包成docker镜像。


### 物理机运行微服务

```bash
$ cd build/
$ ./service #运行service 
2019/11/11 15:04:48 Transport [http] Listening on [::]:45645
2019/11/11 15:04:48 Broker [http] Connected to [::]:41175
2019/11/11 15:04:48 Registry [mdns] Registering node: wida.micro.srv.greeter-ffb7abee-f8e2-43fe-a238-8f8348172dae
2019/11/11 15:07:59 Received Say.Hello request

$ ./cli  #运行cli
&com_hello.SayResponse{Msg:"HelloJohn ,Im wida", Values:[]string(nil), Header:map[string]*com_hello.Pair{"name":(*com_hello.Pair)(0xc0002c0200)}, Type:0, XXX_NoUnkeyedLiteral:struct {}{}, XXX_unrecognized:[]uint8(nil), XXX_sizecache:0} 

$ ./web 
2019/11/11 15:09:54 Listening on [::]:8009
$ curl -d"name=wida" localhost:8009    #这边简化web请求采用curl，同样可以打开浏览器 from表单提交
<html><body><h1>Hello wida ,Im wida</h1></body></html>
```

### 容器化运行服务


## 编写Dockfile

分别编写三个dockerfile
```
FROM alpine:latest

WORKDIR /
COPY  build/cli /

CMD ["./cli"]
```

#### 编写docker-composer.yml

```
version: "3"

services:
  etcd:
    command: --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379
    image: appcelerator/etcd:latest
    ports:
      - "2379:2379"
    networks:
      - overlay
  server:
    command: ./service --registry=etcd --registry_address=etcd:2379 --register_interval=5 --register_ttl=10
    build:
      dockerfile: ./service.Dockerfile
      context: .
    image: wida/micro-service:v1.0
    networks:
      - overlay
    depends_on:
      - etcd
    deploy:
      mode: replicated
      replicas: 2
  web:
    command: ./web --registry=etcd --registry_address=etcd:2379 --register_interval=5 --register_ttl=10
    image: wida/micro-web:v1.0
    build:
      dockerfile: ./web.Dockerfile
      context: .
    networks:
      - overlay
    depends_on:
      - etcd
    ports:
      - "8009:8009"
networks:
  overlay:
```

这个时候我们我们采用了`etcd`做服务发现。`etcd`相关文档请查看[官方文档](https://etcd.io/)

#### docker-compose 部署服务

```bash
$ make docker
$ curl -d "name=wida" http://127.0.0.1:8009 ## 看服务都完全部署完毕后执行
<html><body><h1>Hello 111 wida ,Im 7247ec528ca4</h1></body></html>
```

### 在Kuberbetes中运行
coming soon


## 总结

本小节简要的介绍了go-micro的使用，go-micro实践的生态也比较庞大，代码也写的很优雅是一个不错的代码研究学习对象。更多关于go-micro的情况可以看github源码和参考资料中的文档。

本小节[代码](https://gitlab.ulucu.com/xcxia/learning-go-code/tree/master/mircotest)


# 参考资料

- [Micro中文文档](https://micro.mu)
- [Micro 中国站教程系列](https://github.com/micro-in-cn/tutorials)