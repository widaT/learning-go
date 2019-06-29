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


参数的介绍

- o 指定输出的文件名，可以带上路径，例如 go build -o a/b/c
- i 安装相应的包
- a 更新全部已经是最新的包的，但是对标准包不适用
- n 把需要执行的编译命令打印出来，但是不执行，这样就可以很容易的知道底层是如何运行的
- p n 指定可以并行可运行的编译数目，默认是CPU数目
- race 开启编译的时候自动检测数据竞争的情况，目前只支持64位的机器
- v 打印出来我们正在编译的包名
- work 打印出来编译时候的临时文件夹名称，并且如果已经存在的话就不要删除
- x 打印出来执行的命令，其实就是和-n的结果类似，只是这个会执行
- ccflags 'arg list' 传递参数给5c, 6c, 8c 调用
- compiler name 指定相应的编译器，gccgo还是gc
- gccgoflags 'arg list' 传递参数给gccgo编译连接调用
- gcflags 'arg list' 传递参数给5g, 6g, 8g 调用
- installsuffix suffix 为了和默认的安装包区别开来，采用这个前缀来重新安装那些依赖的包，- race的时候默认已经是-installsuffix race,大家可以通过-n命令来验证
- ldflags 'flag list' 传递参数给5l, 6l, 8l 调用
- tags 'tag list' 设置在编译的时候可以适配的那些tag，详细的tag限制参考里面的 Build Constraints


## go get

这个命令是用来动态获取远程代码包的，目前支持的有BitBucket、GitHub、Google Code和Launchpad。这个命令在内部实际上分成了两步操作：第一步是下载源码包，第二步是执行go install。下载源码包的go工具会自动根据不同的域名调用不同的源码工具，对应关系如下：
BitBucket (Mercurial Git)
GitHub (Git)
Google Code Project Hosting (Git, Mercurial, Subversion)
Launchpad (Bazaar)
所以为了go get 能正常工作，你必须确保安装了合适的源码管理工具，并同时把这些命令加入你的PATH中。其实go get支持自定义域名的功能，具体参见go help remote。
参数介绍：
-d 只下载不安装
-f 只有在你包含了-u参数的时候才有效，不让-u去验证import中的每一个都已经获取了，这对于本地fork的包特别有用
-fix 在获取源码之后先运行fix，然后再去做其他的事情
-t 同时也下载需要为运行测试所需要的包
-u 强制使用网络去更新包和它的依赖包
-v 显示执行的命令

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
