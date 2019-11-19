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

## 类型转换

### 数组类型转换

数值类型对应表

|  c   | cgo  |
| ------------- |:-------------:|
| signed char  | C.char |
| unsigned char  | C.uchar  |
| unsigned short |C.short |
| unsigned int  | C.int |
| unsigned long  | C.long |
|long long |C.longlong |
|unsigned long long |C.ulonglong |
|float |C.float|
|double |C.double |
|void*|unsafe.Pointer|


### 结构体，联合体，枚举类型转换

cgo使用`C.struct_xxx` 访问c语言中的结构体，例如c语言中结构体为`struct S ` cgo返回为`C.struct_S`。

cgo使用`C.union_xxx` 访问c语言中的联合体，注意的是cgo没办法访问联合体内的字段，`C.union_xxx`会变成字节数组。

cgo使用`C.enum_xxx`访问c语言枚举类型。


```go
package main
/*
struct S {
    int i;
    float type;yuy
union U {
    int i;
    float f;
};

enum E {
    A,
    B,
};

*/
import "C"
import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

func main()  {
	var a C.struct_S
	a.i =10
	a._type =10.0
	fmt.Println(a.i)
	fmt.Println(a._type)


	var b C.union_U  //联合体转成了[4]byte

	binary.LittleEndian.PutUint32(b[:],9) //写入i的值

	fmt.Printf("%T\n", b) // [4]uint8
	fmt.Println(*(*C.int)(unsafe.Pointer(&b)) )

	var c C.enum_E = C.B
	fmt.Println(c)
}
```

```bash
$ go run main.go
10
10
[4]uint8
9
1
```

### 字符串和字节数组

go string转 char * ,[]byte 转 void *(unsafe.Pointer)
```go
func C.CString(string) *C.char
func C.CBytes([]byte) unsafe.Pointer
```
这几个转换都有额外的内存耗损，使用完记得`C.free`释放内存。


char * 转go string，void *(unsafe.Pointer)转[]byte
```go
func C.GoString(*C.char) string
func C.GoStringN(*C.char, C.int) string

func C.GoBytes(unsafe.Pointer, C.int) []byte
```


## cgo中c函数返回值

c语言是不支持多个返回值的，cgo调用c函数会有两个返回值（即使c函数原本返回void）。第二个返回值为
`#include <errno.h>`中的errno变量对应的错误描述，errno是一个全局变量，用于返回最近一次调用错误结果。

```go
package main
/*
#cgo LDFLAGS:-lm
#include <math.h>
#include <errno.h>
int div(int a, int b) {
    if(b == 0) {
        errno = EINVAL;
        return 0;
    }
    return a/b;
}
void voidFunc() {}

*/
import "C"
import "fmt"

func main() {
	_, err := C.sqrt(-1) //参数错误
	fmt.Println(err)

	n,err := C.sqrt(4) 
	fmt.Println(n,err)

	_,err = C.voidFunc() //函数返回void
	fmt.Println(err)

	d,err :=C.div(4,0) //除数不能为0
	fmt.Println(d,err)
}
```

```go
$ go run main.go
numerical argument out of domain
2 <nil>
<nil>
0 invalid argument
```

## 参考资料

-   [cgo](https://golang.org/cmd/cgo/)