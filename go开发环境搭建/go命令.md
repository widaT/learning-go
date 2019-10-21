# go 命令

go语言本身自带一个命令行工具，这些命令在项目开发中会反复的使用，我们发点时间先了解下，像在上个小节我们使用的go evn 查看go相关的环境变量，还有使用go build 编译了我们写的hello world。我们看下完整的go 命令：

```
$ go
Go is a tool for managing Go source code.

Usage:
	go <command> [arguments]
The commands are:
	bug         start a bug report
	build       compile packages and dependencies
	clean       remove object files and cached files
	doc         show documentation for package or symbol
	env         print Go environment information
	fix         update packages to use new APIs
	fmt         gofmt (reformat) package sources
	generate    generate Go files by processing source
	get         download and install packages and dependencies
	install     compile and install packages and dependencies
	list        list packages or modules
	mod         module maintenance
	run         compile and run Go program
	test        test packages
	tool        run specified go tool
	version     print Go version
	vet         report likely mistakes in packages
```


在这些go命令中，有些命令相对功能简单，有些则相对复杂些，我们分几个小节介绍这些命令的使用场景。本小节我们先看 go build 和 go get。

## go build


go build命令 编译我们指定的go源码文件，或者指定路径名的包，以及它们的依赖包。
如果我们在执行go build命令时没有任何参数，那么该命令将试图编译当前目录所对应的代码文件。

在编译main 包时，将以第一个参数文件或者路径，作为输出的可执行文件名，例如
第一个源文件（'go build ed.go rx.go' 生产'ed'或'ed.exe'）
或源代码目录（'go build unix/sam'生成'sam'或'sam.exe'）。

编译多个包或单个非main包时，go build 编译包但丢弃生成的对象，这种情况仅用作检查包是否可以顺利编译通过。

go build 编译包时，会忽略以“_test.go”结尾的文件。

我们平常在开发一些兼容操作系统底层api的项目时，我们可以根据相应的操作系统编写不同的兼容代码，例如我们在处理signal的项目，linux和Windows signal是有差异的，我们可以通过加操作系统后缀的方式来命名文件，例如sigal_linux.go sigal_windows.go，在go build的时候会更加当前的操作系统（$GOOS 环境变量）来下选择编译文件，而忽略非该操作系统的文件。上面的例子如果在linux下 go build 会编译 sigal_linux.go  而忽略 sigal_windows.go，而windows系统下在反之。

go build 有好多的个参数，其中比较常用的是 

- -o 指定输出的文件名，可以带上路径，例如 go build -o a/b/c。
- -race 开启编译的时候检测数据竞争。
- -gcflags 这个参数在编译优化和逃逸分析中经常会使用到。

### go build -tags
上面我们介绍了go build 可以使用加操作系统后缀的文件名来选择性编译文件。还有一种方法来实现条件编译，那就是使用`go build -tags tag.list`
我创建一个项目`buildtag-demo`，目录架构为
```
.
├── go.mod
├── main.go
├── tag_d.go
└── tag_p.go
```
我们在`main.go` 里头的代码是
```golang
package main

func main()  {
	debug()
}
```

`tag_d.go` 代表debug场景下要编译的文件，文件内容是
```golang
// +build debug

package main

import "fmt"

func debug()  {
	fmt.Println("debug ")
}
``` 
`tag_p.go` 代表生产环境下要编译的文件，文件内容是

```golang
// +build !debug

package main
func debug()  {
	//fmt.Println("debug ")
}
```

我们编译和运行下项目
```bash
$ go build
$ ./buildtag-demo
//nothing 没有任何输出
$ go build -tags debug
$ ./buildtag-demo
debug
```

#### 使用方式
- `+build`注释需要在`package`语句之前，而且之后还需要一行空行。
- `+build`后面跟一些条件 只有当条件满足的时候才编译此文件。
- `+build`的规则不仅见约束`.go`结尾的文件，还可以约束go的所有源码文件。


#### `+build`条件语法：
- 只能是是数字、字母、_（下划线）
- 多个条件之间，空格表示OR；逗号表示AND；叹号(!)表示NOT
- 可以有多行`+build`，它们之间的关系是AND。如：
	```
	// +build linux darwin
	// +build 386
	等价于
	// +build (linux OR darwin) AND 386
	```
- 加上`// +build ignore`的原件文件可以不被编译。

## go get

go get 在go modules出现之前一直作为go获取依赖包的工具，在go modules出现后，go get的功能和之前有了不一样的定位。现在go get主要的功能是获取解析并将依赖项添加到当前开发模块然后构建并安装它们。

参考 go modules 章节

## go还提供了其它很多的工具，例如下面的这些工具

- go bug 给go官方提交go bug（执行命令后会在浏览器弹出go github issue提交页面）
- go clean 这个命令是用来移除当前源码包和关联源码包里面编译生成的文件.
- go doc  这个命令用来从包文件中提取顶级声明的首行注释以及每个对象的相关注释，并生成相关文档。
- go env 这个命令我们之前介绍过，是用来获取go的环境变量
- go fix 用于将你的go代码从旧的发行版迁移到最新的发行版，它主要负责简单的、重复的、枯燥无味的修改工作，如果像 API 等复杂的函数修改，工具则会给出文件名和代码行数的提示以便让开发人员快速定位并升级代码。
- go fmt 格式化go代码。
- go generate 这个命令用来生成go代码文件，平常工作中比较少接触到。
- go install 编译生产可执行文件，同时将可执行文件移到$GOBIN这个环境变量设置的目录下
- go list 查看当前项目的依赖模块（包），在go modules 小节会看到一些具体用法
- go mod go的包管理工具，go modules会单独介绍
- go run 编译和运行go程序
- go test go的单元测试框架，会在go test小节单独介绍
- go tool 这个命令聚合了很多go工具集，主要关注 go tool pprof 和 go tool cgo这两个命令，在go 性能调优和cgo的章节会讲到这两个命令。
- go version 查看go当前的版本
- go vet 用来分析当前目录的代码是否正确。



# 参考资料

- [go build -tags 试验](https://www.jianshu.com/p/858a0791f618)