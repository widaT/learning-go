# go 1.18 workspace使用

go语言对于模块化代码组织方式说实话，不是很理想。从最早被人诟病的`go path`方式，包括后来稍微有点现代语言模块化方式的`go modules`也槽点满满。

虽然不尽人意，但是go官方都是没有放弃继续改进模块化代码组织方式。这次go1.18又有了新的一个功能叫做 `go workspace`，中文翻译为go工作区。

## 初识 go workspace

### 什么需要 go workspace？

我看了下工作区的[提案](https://go.googlesource.com/proposal/+/master/design/45713-workspace.md)，说了`go workspace`设计的初衷

```
在一个大项目中依赖多个go mod项目，而我们需要同时修改go mod项目代码，在原先go mod设计中，依赖的项目是只读的。
```

`workspace`方式是对现有`go mod`方式的补充而非替换，`workspace`会非常方便的在`go`语言项目组加入本地依赖库。

## 使用 go workspace
`go` 命令行增加了 `go work`来支持`go workspace`操作。

```bash
$ go help work
Go workspace provides access to operations on workspaces.
...
Usage:

        go work <command> [arguments]

The commands are:

        edit        edit go.work from tools or scripts
        init        initialize workspace file
        sync        sync workspace build list to modules
        use         add modules to workspace file
...
```

我们看到 `go work` 有四个子命令.

接下来我们创建几个子项目：
```bash
$ go work init a b c
go: creating workspace file: no go.mod file exists in directory a
```
额，报错了。a,b,c 一定是`go mod`项目才能使用.我们创建a,b,c目录，同时添加go mod
```bash
$ tree .
.
├── a
│   └── go.mod
├── b
│   └── go.mod
├── c
    └── go.mod

```

```bash
$ go work init a b c
$ tree .
.
├── a
│   └── go.mod
├── b
│   └── go.mod
├── c
│   └── go.mod
└── go.work
```
ok,工作区创建好了。我们规划下，a目录项目编译可执行文件，b，c为lib库。在工作区目录下`a，b,c 三个模块互相可见`。

```bash
$ tree 
.
├── a
│   ├── a.go
│   └── go.mod
├── b
│   ├── b.go
│   └── go.mod
├── c
│   ├── c.go
│   └── go.mod
└── go.work
```

b.go

```go
import "fmt"

func B() {
	fmt.Println("I'm mod b")
}

```

c.go

```go
package c

import (
	"b"
	"fmt"
)

func C() {
	b.B()
	fmt.Println("I'm mod c")
}

```

a.go

```go
package main

import (
	"b"
	"c"
)
func main() {
	b.B()
	c.C()
}
```
正如你所看到的，a 依赖 b，c依赖b。好了，我们run一下程序

```bash
$ go run a/a.go
I'm mod b
I'm mod b
I'm mod c
$ cd a/    #切换到 a目录下看是不是可以
$ go run a.go # a目录下 b，c也对a可见
I'm mod b
I'm mod b
I'm mod c
```

## 关于是否提交go.work文件

我们看到很多文章说不建议提交go.work文件，说实话我看到这个建议很奇怪，rust项目管理工具`cargo`也有类似的工作区概念，cargo的项目肯定会提交工作区文件，因为这个文件本身是项目的一部分。很多文章说`go.work`这个文件主要用于本地开发，我倒觉得未必有啊，一个多模块，大团队项目不是也可以以工作区项目开发吗？工作区可以非常清楚的划分功能模块，在一个代码仓库里头没啥问题吧。可能这一块会有很多争议，目前用的人少，等大规模运用了在看看情况。

