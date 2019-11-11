# go-micro微服务框架

## go-micro简介
Go Micro是可插拔的微服务开发框架。micro的设计哲学是可插拔的架构理念，她提供可快速构建系统的组件，并且可以根据自身的需求剥离默认实现并自行定制。详细的介绍可参考官方中文文档[Go Micro](https://micro.mu/docs/cn/go-micro.html)

## 安装go-micro依赖
    
-   安装protobuf 参考[grcp和protobuf](./grpc和protobuf.md) 中protobuf的安装方式
-   安装protoc-gen-micro，这个是go-micro定制的protobuf插件，用于生成go-micro定制的类似grpc的代码。
    ```
    go get github.com/micro/protoc-gen-micro
    ```
## 写一个微服务版本的hello world程序

这个程序包含一个service，web，cli：
- service作为微服务的具体体现，他的功能是提供某一类服务，我们这边的稍作简化，只提供一个grpc服务。在实践场景中service是某一类服务接口的集合，接口数量不宜太多，得考虑负载情况和业务情况相结合做适当的拆分整合。
- web作为前端应用，对外提供http服务。在实践场景中web通常会整合各种后端service包装成业务接口提供给终端或者第三方合作伙伴使用。
- cli也是前端应用的另一种形式，不对外提供http服务，可以在命令行下的脚本任务或者某些守护进程。

### 跟grpc一样我们先编写proto文件，定义service

```protobuf
syntax = "proto3";
package com.hello;
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

我们写一个shell程序来生成golang文件
```bash
#!/bin/bash
protoc --go_out=plugins=micro:. test.proto
```
目录结构
```bash
$ tree
.
├── cmd
│   ├── cli
│   │   └── main.go
│   ├── service
│   │   └── main.go
│   └── web
│       └── main.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── proto
│   ├── gen.sh
│   ├── test.pb.go
│   └── test.proto
└── README.md
```
### 编写service
注意go-micro这边使用了默认的mdns来做服务发现，mdns的相关原理可参考[Multicast_DNS](https://en.wikipedia.org/wiki/Multicast_DNS)

```golang
package main

import (
	"context"
	"github.com/micro/go-micro"
	"log"
	hello "mircotest/proto"
	"os"
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
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/web"
	"log"
	hello "mircotest/proto"
	"net/http"
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
			cl := hello.NewSayClient("wida.micro.srv.greeter", client.DefaultClient)
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
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/wrapper/select/roundrobin"
	hello "mircotest/proto"
)

func main() {
	wrapper := roundrobin.NewClientWrapper()
	service := micro.NewService(
		micro.WrapClient(wrapper),
	)
	service.Init()
	cl := hello.NewSayClient("wida.micro.srv.greeter", service.Client())
	rsp, err := cl.Hello(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v \n", rsp)
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
    @mkdir -p build
	$(GOBUILD) -a -installsuffix cgo -ldflags '-w' -o build/web cmd/web/*.go

cli :
	@echo "build cli"
    @mkdir -p build
	$(GOBUILD) -a -installsuffix cgo -ldflags '-w' -o build/cli cmd/cli/*.go

.PHONY: clean
clean:
	@rm -rf build/

.PHONY: proto
proto:
	protoc --go_out=plugins=micro:. proto/test.proto
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
│   ├── cli
│   ├── service
│   └── web
├── cmd
│   ├── cli
│   │   └── main.go
│   ├── service
│   │   └── main.go
│   └── web
│       └── main.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── proto
│   ├── gen.sh
│   ├── test.pb.go
│   └── test.proto
└── README.md

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

####  制作镜像

编写Dockerfile

```
FROM alpine:latest

WORKDIR /
COPY  cmd/service/service /

#EXPOSE 8009       # web的时候去掉注释
CMD ["./service"]
```

#### 生成镜像

``` bash
$ docker build -t wida/micro-service:v1.0 .
Sending build context to Docker daemon  68.66MB
Step 1/4 : FROM alpine:latest
 ---> 4d90542f0623
Step 2/4 : WORKDIR /
 ---> Using cache
 ---> 48b2993f945d
Step 3/4 : COPY  build/service /
 ---> Using cache
 ---> b90abc27c4bc
Step 4/4 : CMD ["./service"]
 ---> Using cache
 ---> b1dcd224c140
Successfully built b1dcd224c140
Successfully tagged wida/micro-service:v1.0
$ docker images
REPOSITORY            TAG                 IMAGE ID            CREATED             SIZE
wida/micro-service          v1.0                 b1dcd224c140        9 seconds ago       28.3MB
```

使用同样的方式生成web

```bash
docker images wida/*
REPOSITORY           TAG                 IMAGE ID            CREATED             SIZE
wida/micro-service   v1.0                b1dcd224c140        5 minutes ago       28.3MB
wida/micro-web       v1.0                697738bd1ff0        5 minutes ago       28.4MB
```

#### 编写docker-composer.yml

```
version: "3"
services:
  consul:
    command: -server -bootstrap -rejoin
    image: progrium/consul:latest
    ports:
      - "8300:8300"
      - "8400:8400"
      - "8900:8500"
      - "8600:53/udp"
    networks:
      - overlay
  server:
    command: ./service --registry=consul --registry_address=consul:8500 --register_interval=5 --register_ttl=10
    image: wida/micro-service:v1.0
    networks:
      - overlay
    deploy:
      mode: replicated
      replicas: 2
  web:
    command: ./web --registry=consul --registry_address=consul:8500 --register_interval=5 --register_ttl=10
    image: wida/micro-web:v1.0
    networks:
      - overlay
    ports:
      - "8009:8009"
networks:
  overlay:
```

这个时候我们我们采用了`consul`做服务发现。`cunsul`相关文档请查看[官方文档](https://www.consul.io/docs/index.html)

#### docker swarm部署服务

```bash
$ docker stack deploy -c docker-compose.yml go-micro
Creating network go-micro_overlay
Creating service go-micro_web
Creating service go-micro_consul
Creating service go-micro_service
$ docker service ls # 查看下服务状态 主要看replicas
ID                  NAME                MODE                REPLICAS            IMAGE                    PORTS
i6glmed3o9hs        go-micro_consul     replicated          1/1                 progrium/consul:latest   *:8300->8300/tcp, *:8400->8400/tcp, *:8900->8500/tcp, *:8600->53/udp
oljgulgk4dby        go-micro_service     replicated          2/2                 wida/micro-service:v1.0   
qki1bb7wdog2        go-micro_web        replicated          1/1                 wida/micro-web:v1.0      *:8009->8009/tcp
$ curl -d "name=wida" http://127.0.0.1:8009 ## 看服务都完全部署完毕后执行
<html><body><h1>Hello 111 wida ,Im 7247ec528ca4</h1></body></html>
$ curl -d "name=wida" http://127.0.0.1:8009  
<html><body><h1>Hello 111 wida ,Im 31c430aa0aaf</h1></body></html> # 注意这边的hostname 和上面那个请求不一样，说明负载均衡生效了
```

### 在Kuberbetes中运行
coming soon



## 总结

本小节简要的介绍了go-micro的使用，go-micro实践的生态也比较庞大，代码也写的很优雅是一个不错的代码研究学习对象。更多关于go-micro的情况可以看github源码和参考资料中的文档


# 参考资料

- [Micro中文文档](https://micro.mu/docs/cn/)
- [Micro 中国站教程系列](https://github.com/micro-in-cn/tutorials)