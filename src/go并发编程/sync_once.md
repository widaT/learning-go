# sync.Once

sync.Once提供一种机制来保证某个函数只执行一次，常常用于初始化对象。
只有一个api：
```go
func (o *Once) Do(f func())
```

我们来写个demo

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	once :=	sync.Once{}
	a :=1
	i:=0
	for i <10{
		go func() {
			once.Do(func() {
				a++
			})
		}()
		i++
	}
	time.Sleep(1e9)
	fmt.Println(a)
}
```

```bash
$ go run main.go
2
```
上面的例子中，Once.Do 被执行了十次，他包裹的函数却只被执行了一次。

我们来看下`sync.Once`的源码：

```go
type Once struct {
	done uint32  //状态值
	m    Mutex  //互斥锁
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 { //等于0的时候代表没被执行过
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()              //确保只有一个goroutine能到锁
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)  //原子操作改变状态
		f()
	}
}
```

从源码上看，他是利用`sync.Mutex`和原子操作来实现的。
