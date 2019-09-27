# go channal

我们先通过golang汇编追踪下 make（chan，int）到底做了什么事情。
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

那么`ch := make(chan int,1)` 其实是调用 runtime.makechan的方法，返回的类型是*hchan。


`ret := <- ch` 这个语句对于汇编是
```bash
LEAQ    ""..autotmp_12+64(SP), CX
MOVQ    CX, 8(SP)
CALL    runtime.chanrecv2(SB)
```

`ch <-1`对应的汇编是
```
LEAQ    ""..stmp_0(SB), CX
MOVQ    CX, 8(SP)
CALL    runtime.chansend1(SB)
```
我看从go源码的（src/runtime/chan.go）看到的hchan struct
```golang
type hchan struct {
	qcount   uint           // 当前使用的个数
	dataqsiz uint           // 缓冲区大小 make的第二个参数
	buf      unsafe.Pointer // 缓存数组的指针
	elemsize uint16         // 元素size
	closed   uint32         //是否已经关闭 1已经关闭 0还没关闭
	elemtype *_type // 元素类型
	sendx    uint   // 发送的索引 比如缓冲区大小为3 这个索引 经历 0-1-2 然后再从0开始
	recvx    uint   // 接收的索引 
	recvq    waitq  // 等待recv的goroutine链表
	sendq    waitq  // 等待send的goroutine链表
	lock mutex      //互斥锁
}

type waitq struct {
	first *sudog
	last  *sudog
}

type sudog struct {
	g *g
	isSelect bool
	next     *sudog
	prev     *sudog
	elem     unsafe.Pointer // data element (may point to stack)
	acquiretime int64
	releasetime int64
	ticket      uint32
	parent      *sudog // semaRoot binary tree
	waitlink    *sudog // g.waiting list or semaRoot
	waittail    *sudog // semaRoot
	c           *hchan // channel
}
```


获取chan的关闭状态

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