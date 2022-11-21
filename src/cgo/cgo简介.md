# cgo简介

## 什么是cgo

在某些场景下go很可能需要调动c函数，比如调用系统底层驱动。c在过去几十年的积累了非常多优秀的lib库，所以很多编程语言都会选择跨语言调用c，这项技术被称为foreign-function interfaces （简称ffi）。go和c语言有相当深的颜渊，go当然不会拒绝继承c语言的这些财产，于是便有了cgo。

cgo是golang自带的go和c互相调用的工具。注意go不仅仅可以调用c，而且还可以让c调用go。

## cgo is not go

任何技术都不是完美的，虽然go可以通过cgo的方式使用c语言的很多优秀的lib库，但是cgo的使用并不是没有代价的，在某些场景下代价还挺大。[cgo is not go](https://dave.cheney.net/2016/01/18/cgo-is-not-go)这个文章详细描述了cgo不是很让人满意的地方。

总的来说有如下几点：

- cgo让构建go程序变得复杂，必须安装c语言相关工具链
- cgo不支持交叉编译
- cgo无法使用go语言很多工具，像pprof
- 用了cgo后部署变得更加复杂。c语言很多依赖动态链接库，go语言原先只编成一个二进制可执行文件，现在这种局面被打破。
- cgo代理额外的性能损耗。cgo中go的堆栈和c语言堆栈是分开的，go调用c语言函数通常需要申请额外的堆内存进行数据拷贝。

cgo作用大，但是非必要条件不要使用cgo。

## cgo hello world

```go
package main

/*
#include <stdio.h>
#include <stdlib.h> //C.free 依赖这个头文件

void myprint(char* s) {
	printf("%s\n", s);
}
*/
import "C"

import "unsafe"

func main() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}
```

```bash
$ go run main.go
Hello from stdio
````

## 参考资料

[cgo-is-not-go](https://dave.cheney.net/2016/01/18/cgo-is-not-go)