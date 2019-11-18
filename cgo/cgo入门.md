# cgo入门

## 开启cgo

要使用cgo需要先导入伪Package `C`，

```golang
// #include <stdio.h>     //这边的特殊注释属于c语言的代码
// #include <errno.h>
import "C" //这边需要注意 上面不能有空行，而且这行不能和其他import同行
```
当然还可以换一种注释

```golang
/* 
#include <stdio.h>
#include <errno.h>
*/
import "C" 
```

import `C`上面注释里头可以写c语言的代码

```golang
/* 
#include <stdio.h>
#include <errno.h>
void sayHello() {
    printf("hello");
}
*/
import "C" 
```

可以引入自定义的头文件

```bash
$  tree
.
├── go.mod
├── main.go
├── myhead.h
└── say.c
```
myhead.h内容

```c
#include<stdio.h>
#include<stdlib.h>

int sayHello(void * a, const char* b);
```

say.c内容

```c
#include "myhead.h"

int sayHello(void *a, const char* b){
	return sprintf((char *)a,b);
}
```

main.go 内容
```golang
package main

/*
#include "myhead.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)
func main() {
	b := make([]byte,20)
	bb := C.CBytes(b)         //会拷贝内存
	char :=C.CString("Hello, World") //会拷贝内存
	defer func() {    //这边需要手动释放堆内存
		C.free(unsafe.Pointer(char))
		C.free(bb)
	}()
	l := C.sayHello(bb,char)
	fmt.Println(string(C.GoBytes(bb,l)))
}
```
编译的时候go编译器会自动扫描源码目录下c语言的源码文件
```bash 
$ go build  #和正常程序编译没有区别
$ tree
.
├── cdo-demo    #这个是编译后的程序
├── go.mod
├── main.go
├── myhead.h
└── say.c
$ ./cdo-demo
Hello, World
```

## cgo的编译链接参数

cgo的编译链接参数通过在注释中使用 `#cgo`伪命令实现。`#cgo`命令定义CFLAGS，CPPFLAGS，CXXFLAGS，FFLAGS和LDFLAGS，以调整C，C ++或Fortran编译器的行为。多个指令中定义的值被串联在一起。同时`#cgo`可以包括一系列构建约束，这些约束将其作用限制在满足其中一个约束的系统上。

```go
// #cgo CFLAGS：-DPNG_DEBUG = 1 #宏定义及赋值
// #cgo amd64 386 CFLAGS：-DX86 = 1  #这边包好了系统约束 （adm64或者386）
// #cgo LDFLAGS：-lpng 
// #include <png.h> 
import“ C”
```

另外，CPPFLAGS和LDFLAGS可以通过pkg-config工具使用`#cgo pkg-config:`指令后跟软件包名称来获得。例如：
```go
// #cgo pkg-config: png cairo
// #include <png.h>
import "C
```
可以通过设置PKG_CONFIG环境变量来更改默认的pkg-config工具。

在程序构建时：
- 程序包中的所有CPPFLAGS和CFLAGS伪指令被串联在一起，并用于编译该软件包中的C文件。
- 程序包中的所有CPPFLAGS和CXXFLAGS伪指令被串联在一起，并用于编译该程序包中的C++文件。
- 程序包中的所有CPPFLAGS和FFLAGS指令都已连接在一起，并用于编译该程序包中的Fortran文件。
- 程序中任何程序包中的所有LDFLAGS伪指令在链接时被串联并使用。
- 所有pkg-config伪指令被串联并同时发送到pkg-config，以添加到适当的编译链接参数。


`#cgo`伪命令中包含`${SRCDIR}`标志的将被展开成源文件的目录的绝对路径。
```go
// #cgo LDFLAGS: -L${SRCDIR}/libs -lfoo
```
上面的代码在将被展开为：
```go
// #cgo LDFLAGS: -L/go/src/foo/libs -lfoo
```


## 参考资料

-   [cgo](https://golang.org/cmd/cgo/)