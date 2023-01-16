# 第一个go程序

go环境已经搭建好了，接下来我们写一下go程序的`hello world`


## Hello World

go使用go mod 来管理依赖，所以我们先创建一个go语言项目。

```bash
$ mkdir helloworld && cd helloworld #创建目录
$ go mod init helloworld       #go mod初始化项目，项目名为helloworld
$ touch main.go                #创建代码文件
```

编辑 main.go
```go
package main
import (
    "fmt" //导入fmt package
)

func main() {
    fmt.Println("hello world")
}
```

使用 `go run` 临时编译执行go程序

```bash
$ go run main.go
hello world
```

到此，我们的Go环境基本搭建好了，接下来我们学习过程中需要不断的写练习代码。