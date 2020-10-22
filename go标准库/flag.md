# flag包 —— 解析命令行参数

我们通过`io`包，知道golang使用`os.Args`来接收命令行参数，在简单的场景下`os.Args`基本够我们使用。
复杂的参数场下golang提供`flag`包来解析命令行参数。


## 参数格式

```
-flag　　//只支持布尔类型，有为true，无为默认值
-flag=x
-flag x  // bool型不能用这个方式
```

参数中带`-`和`--`功能一致，例如`cmd -a 3`和`cmd --a 3`是一致的。

## 定义flag接收参数的两种方式

1. 接收指针

例如 `flag.String(), flag.Bool(), flag.Int()`

```go
func Int(name string, value int, usage string) *int
func Bool(name string, value bool, usage string) *bool
func String(name string, value string, usage string) *string
```
这类函数，`name`是参数名，`value`为默认值，`usage`为帮助信息`-h`的时候会打印，返回值为指针类型。

```go
import "flag"
var nFlag = flag.Int("n", 1234, "help message for flag n")

fmt.Println(*nFlag) //指针类型用使用需要 *符号
```

2. 参数绑定

例如`flag.IntVar()，flag.StringVar(),flag.BoolVar()`

```go
func IntVar(p *int, name string, value int, usage string)
func BoolVar(p *bool, name string, value bool, usage string)
func StringVar(p *string, name string, value string, usage string)
```

这类函数，p为需要的绑定的指针，`name`是参数名，`value`为默认值，`usage`为帮助信息`-h`的时候会打印。

```go
var flagvar int
flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
fmt.Println(flagvar)
```

## flag.Parse()

所以的参数能够解析一定要调用`flag.Parse()`，`flag.Parse()`调用后才开始做参数解析。


## 例子

```go
package main

import (
	"flag"
	"fmt"
)

func main() {

	var a int
	var b bool
	var c string

	flag.IntVar(&a, "i", 0, "int flag value")
	flag.BoolVar(&b, "b", false, "bool flag value")
	flag.StringVar(&c, "s", "default", "string flag value")
	flag.Parse()

	fmt.Println("a:", a)
	fmt.Println("b:", b)
	fmt.Println("c:", c)
}
```

```bash
$ ./cmd -h    #打印帮助信息
Usage of ./cmd:
  -b    bool flag value
  -i int
        int flag value
  -s string
        string flag value (default "default")

$ ./cmd -i 1 -b -s="acb"
a: 1
b: true
c: acb
```


# 参考资料

[go pkg flag](https://golang.org/pkg/flag/)