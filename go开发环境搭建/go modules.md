# go modules

## go modules 介绍

go modules是go 1.11中的一个实验性选择加入功能。目前随着go 1.12 的发布，越来越多的项目已经采用go modules方式作为项目的包依赖管理。

### 设置 GO111MODULE

GO111MODULE 有三个值 off, on和auto（golang 1.12默认值）。

- off：go tool chain 不会支持go module功能，寻找依赖包的方式将会沿用旧版本那种通过vendor目录或者GOPATH/src模式来查找。
- on：go tool chain 会使用go modules，而不会去GOPATH/src目录下查找,依赖包文件保持在$GOPATH/pkg下，允许同一个package多个版本并存，且多个项目可以共享缓存的module。
- auto（go 1.12的默认值）：go tool chain将会根据当前目录来决定是否启用module功能。当前目录在$GOPATH/src之外且该目录包含go.mod文件或者当前文件在包含go.mod文件的目录下面则会开启 go modules。


### go mod 命令介绍

go modules 在golang中使用go mod命令来实现。

```
$ go help mod 

download    download modules to local cache (下载依赖包到本地)
edit        edit go.mod from tools or scripts　（编辑go.mod）
graph       print module requirement graph　（列出模块依赖图）
init        initialize new module in current directory (在当前目录下初始化go module)
tidy        add missing and remove unused modules (下载缺失模块和移除不需要的模块)
vendor      make vendored copy of dependencies （将依赖模块拷贝到vendor下）
verify      verify dependencies have expected content （验证依赖模块）
why         explain why packages or modules are needed （解释为什么需要依赖）
```

## go mod 使用

创建目录
```
$ mkdir -p ~/gocode/hello
$ cd ~/gocode/hello/
````

初始化新module
```
$ go mod init github.com/youname/hello

go: creating new go.mod: module github.com/youname/hello
```
写代码

```
$ cat <<EOF > hello.go
package main

import (
    "fmt"
    "rsc.io/quote"
)

func main() {
    fmt.Println(quote.Hello())
}
EOF
````
构建项目

```
$ go build
$ ./hello
你好，世界。
```

这个时候看下 go.mod 文件 

```
$ cat go.mod

module github.com/you/hello

require rsc.io/quote v1.5.2
```

一旦项目使用go modules的方式解决包依赖，你日常的工作就是在你的go项目代码中添加improt语句，标准命令（go build 或者go test）将根据需要自动添加新的依赖包（更新go.mod并下载新的依赖）。
在需要时可以使用go get foo@v1.2.3，go get foo@master，go get foo@e3702bed2或直接编辑go.mod等命令选择更具体的依赖项版本。


以下还有些常用功能：
```
go list -m all  - 查看将在构建中用于所有直接和间接依赖模块的版本信息
go list -u -m all  - 查看所有直接和间接依赖模块的可用版本升级信息
go get -u     -会升级到最新的次要版本或者修订版本(x.y.z, z是修订版本号， y是次要版本号)
go get -u=patch  - 升级到最新的修订版本
go build ./... 或者go test ./...  - 从模块根目录运行时构建或测试模块中的所有依赖模块
```


### go.mod 文件介绍

go.mod 有module, require、replace和exclude 四个指令
module  指令声明其标识，该指令提供模块路径。模块中所有软件包的导入路径共享模块路径作为公共前缀。模块路径和从go.mod程序包目录的相对路径共同决定了程序包的导入路径。
require 指令指定的依赖模块
replace 指令可以替换依赖模块
exclude 指令可以忽略依赖项模块


### 国内goher使用go modules

由于墙的原因，国内开发者在使用go modules的时候会遇到很多依赖包下载不了的问题，下面提供几个解决方案

- 使用go proxy（推荐） 
```
    export GOPROXY="https://goproxy.io"
```
- 使用go mod replace 替换下载不了的包，例如golang.org下面的包

```
    replace (
        golang.org/x/crypto v0.0.0-20190313024323-a1f597ede03a => github.com/golang/crypto v0.0.0-20190313024323-a1f597ede03a
    )
```



## 补充资料

### 模块版本定义规则

模块必须根据[semver](https://semver.org/)进行语义版本化，通常采用v（major）.（minor）.（patch）的形式，例如v0.1.0，v1.2.3或v1.5.0-rc.1。版本必需是v字母开头。 go mod 在拉取对应包版本的时候会找相应包的git tag（release tag和 pre release tag），如果对应包没有git tag 就会来取master 而对应的版本好会变成v0.0.0。go mod工具中的版本号格式为版本号 + 时间戳 + hash以下的版本都是合法的：

```
gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
github.com/PuerkitoBio/goquery v1.4.1
gopkg.in/yaml.v2 <=v2.2.1
golang.org/x/net v0.0.0-20190125091013-d26f9f9a57f3
latest
```


# 总结

本小节介绍golang modules的使用，go modules本身有很多比较复杂的设计，你可以通过go modules[官方英文文档](https://github.com/golang/go/wiki/Modules#quick-start)做详细了解。go modules对golang 项目构建的基石，在实际项目中一定会经常接触到。 
