# go modules

## go modules 介绍

go 1.11包括对此处提出的版本化模块的初步支持。modules是go 1.11中的一个实验性选择加入功能，其计划是结合反馈并最终确定go 1.13的功能。目前随着go 1.12 的发布，越来越多的项目已经采用go modules方式作为项目的包依赖管理。

### 设置GO111MODULE

GO111MODULE 有三个值 off, on和auto（golang 1.12默认值）。

- off：go tool chain 不会支持go module功能，寻找依赖包的方式将会沿用旧版本那种通过vendor目录或者GOPATH/src模式来查找。
- on：go tool chain 会使用go modules，而不会去GOPATH/src目录下查找,依赖包文件保持在$GOPATH/pkg下，允许同一个package多个版本并存，且多个项目可以共享缓存的module。
- auto（go 1.12的默认值）：go tool chain将会根据当前目录来决定是否启用module功能。当前目录在$GOPATH/src之外且该目录包含go.mod文件或者当前文件在包含go.mod文件的目录下面则会开启 go modules。


### go mod使用

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

## demo

创建目录
```
$ mkdir -p ~/gocode/hello
$ cd ~/gocode/hello/
````

初始化新go module
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

一旦项目使用go modules的方式解决依赖，你日常的工作就是在你的go项目代码中添加improt语句，标准命令（go buidl 或者go test）将根据需要，自动添加新的依赖项（更新go.mod并下载新的依赖项）。
在需要时，可以使用go get foo@v1.2.3，go get foo@master，go get foo@e3702bed2或直接编辑go.mod等命令选择更具体的依赖项版本。


一下还有写常用功能：
```
go list -m all  - 查看将在构建中用于所有直接和间接依赖关系的版本信息
go list -u -m all  - 查看所有直接和间接依赖项的可用版本升级信息
go get -u 或 go get -u=patch  - 更新所有直接和间接依赖关系到最新的次要或补丁升级
go build ./... 或者go test ./...  - 从模块根目录运行时构建或测试模块中的所有软件包
```


### go.mod 文件介绍

go.mod 有module, require、replace和exclude 四个指令
module  语句指定包的名字（路径）
require 语句指定的依赖项模块
replace 语句可以替换依赖项模块
exclude 语句可以忽略依赖项模块





#   参考资料
- [go modules官方文档](https://github.com/golang/go/wiki/Modules#quick-start)