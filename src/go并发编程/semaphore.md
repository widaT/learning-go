# Semaphore（信号量）

## 什么是Semaphore（信号量，信号标）

Semaphore是一种同步对象，用来存0到指定值的计数值。信号量如果计数值可以是任意整数则称为计数信号量（一般信号量），如果值只能是0和1则叫二进制信号量，也就是我们常见的互斥锁（Mutex）。

计数信号量具备两种操作动作，称为V与P。V操作会增加信号量计数器的数值，P操作会减少它。go标准库没有实现信号量，但是拓展同步库里头实现了（`golang.org/x/sync/semaphore`）计数信号量叫做 `Weighted`——加权信号量。

我们先看下`Weighted`结构体
```go
type Weighted struct {
	size    int64
	cur     int64
	mu      sync.Mutex
	waiters list.List
}
```

`Weighted`有四个对外api
```go
func NewWeighted(n int64) *Weighted
func (s *Weighted) Acquire(ctx context.Context, n int64) error
func (s *Weighted) Release(n int64)
func (s *Weighted) TryAcquire(n int64) bool
```

运作方式：

初始化，给与它一个整数值。
运行P（Acquire），信号标S(size)的值将被减少。企图进入临界区段的goroutine，需要先运行P（Acquire）。当信号标S(size)减为负值时，goroutine会被阻塞，不能继续；当信号标S(size)不为负值时，goroutine可以获准进入临界区段。
运行V（Release），信号标S(size)的值会被增加。结束离开临界区段的goroutine，将会运行V（Release）。当信号标S(size)不为负值时，先前被挡住的其他goroutine，将可获准进入临界区段。

TryAcquire的作用和Acquire类似，但是TryAcquire不会阻塞，发现进不去临界区段就会返回false。

## 使用Semaphore

我们在`channel`做goroutine同步的时候介绍了用`channel`控制goroutine数量并发数量的例子。信号量非常适合做类似的事情。
我们写一个控制goroutine并发数为机器cpu核数的程序程序。

```go
package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"
	"golang.org/x/sync/semaphore"
)

func main() {
	ctx := context.TODO()
	var (
		maxWorkers = runtime.GOMAXPROCS(0)
		sem        = semaphore.NewWeighted(int64(maxWorkers))
		out        = make([]int, 8)
	)
	fmt.Println(maxWorkers)
	for i := range out {
		//fmt.Println(sem.TryAcquire(1))
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}
		go func(i int) {
			fmt.Printf("goroutine %d \n",i)
			defer sem.Release(1)
			time.Sleep(1e9)
		}(i)
	}

	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil { //这边会等到size为初始值再返回，这总方式可以实现类似`sycn.WaitGroup`的功能
		log.Printf("Failed to acquire semaphore: %v", err)
	}
}
```

```bash
$ go run main.go
4
goroutine 3 
goroutine 0 
goroutine 1 
goroutine 2 
goroutine 4 #时间间隔1s后打印
goroutine 5 
goroutine 6 
goroutine 7 
```

`Weighted`的`Release`方法会按照FIFO（队列）的顺序唤醒阻塞goroutine。在实际开发中控制并最大并发数的同时还需要防止超长时间的goroutine，所以`sem.Acquire`带了`context`参数。

## 总结

本小节介绍了信号量的运行方式，以及介绍了go semaphore（golang.org/x/sync/semaphore）的使用。