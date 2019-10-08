# go channal是如何实现的

我们先通过golang汇编追踪下 make（chan，int）到底做了什么。

```golang
package  main
import "fmt"
func main()  {
	ch := make(chan int,1)
	go func() {
		for {
			ch <-1
		}
	}()
	ret := <- ch
	fmt.Println(ret)
}
```

```bash
go build --gcflags="-S" main.go
```

看到如下这段go汇编代码
```
 LEAQ    type.chan int(SB), AX //获取type.chan int 类型指针
 MOVQ    AX, (SP)     //压栈  make的第一个参数
 MOVQ    $1, 8(SP)    //压栈  make的第一个参数
 CALL    runtime.makechan(SB)  //实际调用 runtime.makechan函数
````

那么`ch := make(chan int,1)` 其实是调用 runtime.makechan的方法，

go 1.13 源码

```golang
func makechan(t *chantype, size int) *hchan {
	elem := t.elem
	if elem.size >= 1<<16 {
		throw("makechan: invalid channel element type")
	}
	if hchanSize%maxAlign != 0 || elem.align > maxAlign {
		throw("makechan: bad alignment")
	}
	mem, overflow := math.MulUintptr(elem.size, uintptr(size))
	if overflow || mem > maxAlloc-hchanSize || size < 0 {
		panic(plainError("makechan: size out of range"))
	}
	var c *hchan
	switch {
	case mem == 0:
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		c.buf = c.raceaddr()
	case elem.ptrdata == 0:
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		c = new(hchan)
		c.buf = mallocgc(mem, elemblockevent, true)
	}
	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)
	return c
}
```

返回的类型是*hchan。所以chan是一个指针类型，这样子chan在各个goroutine的传递
都是直接传chan，传递chan 指针。

我看从go源码的（src/runtime/chan.go）看到的hchan struct
```golang
type hchan struct {
	qcount   uint           // 当前使用的个数
	dataqsiz uint           // 缓冲区大小 make的第二个参数
	buf      unsafe.Pointer // 缓存数组的指针,带缓冲区的chan指向缓冲区，无缓存的channel指向自己指针（仅做数据竞争分析使用）
	elemsize uint16         // 元素size
	closed   uint32         //是否已经关闭 1已经关闭 0还没关闭
	elemtype *_type // 元素类型
	sendx    uint   // 发送的索引 比如缓冲区大小为3 这个索引 经历 0-1-2 然后再从0开始
	recvx    uint   // 接收的索引 
	recvq    waitq  // 等待recv的goroutine链表
	sendq    waitq  // 等待send的goroutine链表
	lock mutex      //互斥锁
}
```

`ret := <- ch` 这个语句对于汇编是
```bash
LEAQ    ""..autotmp_12+64(SP), CX
MOVQ    CX, 8(SP)
CALL    runtime.chanrecv2(SB)
```

我们看到源码中用 `chanrecv1` 和 `chanrecv2` 两个函数：
```golang
// entry points for <- c from compiled code
//go:nosplit
func chanrecv1(c *hchan, elem unsafe.Pointer) {
	chanrecv(c, elem, true)
}

//go:nosplit
func chanrecv2(c *hchan, elem unsafe.Pointer) (received bool) {
	_, received = chanrecv(c, elem, true)
	return
}
```
我们这边对应的是 `chanrecv2`， 那么什么时候是`chanrecv1`？什么时候是`chanrecv2`呢？
当  `<- ch` chan左边没有接收值的时候使用的是 `chanrecv1`，当左边有接收值的时候是`chanrecv2`，注意 chan 的接收值有两个
```go
ret := <- ch
//带bool值
ret，ok := <- ch
```
这两个都是使用`chanrecv2`。



`ch <-1`对应的汇编是
```
LEAQ    ""..stmp_0(SB), CX
MOVQ    CX, 8(SP)
CALL    runtime.chansend1(SB)
```
这边可以看到实际上是调用`runtime.chansend1` 源码如下
```golang
func chansend1(c *hchan, elem unsafe.Pointer) {
	chansend(c, elem, true, getcallerpc())
}
```
实现上是调用`runtime.chansend`



## 获取chan的关闭状态

通过什么的一些源码阅读我们来实现下chan关闭状态的获取。

```go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	c := make(chan int, 10)
	fmt.Println(isClosed(&c))
	close(c)
	fmt.Println(isClosed(&c))
}

type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
}

func isClosed(c *chan int) uint32 {
	a :=  unsafe.Pointer(c)
	d := (**hchan)(a)
	return (*d).closed
}
```

```bash
# go run main.go
0
1
```

# 参考文档
[深度解密Go语言之channel](https://www.cnblogs.com/qcrao-2018/p/11220651.html#%E6%8E%A5%E6%94%B6)

[understanding channels](https://speakerd.s3.amazonaws.com/presentations/10ac0b1d76a6463aa98ad6a9dec917a7/GopherCon_v10.0.pdf)