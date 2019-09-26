# go channal

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