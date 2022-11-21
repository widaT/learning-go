# goroutine同步

## 互斥锁（sync.Mutex）和读写锁（sync.RWMutex）

类似其他语言，golang也提供了互斥锁和读写锁的的同步原语。

我们先看下go中互斥锁（sync.Mutex）的使用。

```golang
package main
import (
	"fmt"
	"sync"
	"time"
)
func main()  {
	a := 0
	for i:=0;i< 100;i++ {
		go func() {
			a += 1
		}()
	}
	time.Sleep(1e9)
	fmt.Println(a)
	a = 0
	mutex := sync.Mutex{}
	for i:=0;i< 100;i++ {
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			a ++
		}()
	}
	time.Sleep(1e9)
	fmt.Println(a)
}
```

```bash
$ go run main.go
85
100
```

我们在看下读写锁（sync.RWMutex）。

读写锁：读的时候不会阻塞读，会阻塞写；写的时候会阻塞写和读，我们可以用这个特性实现线程安全的map。

```golang
import "sync"
type safeMap struct {
	rwmut sync.RWMutex
	Map map[string]int
}

func (sm *safeMap)Read(key string)(int,bool){
	sm.rwmut.RLock()
	defer sm.rwmut.RUnlock()
	if v,found := sm.Map[key];!found {
		return 0,false
	}else {
		return v,true
	}
}

func (sm *safeMap)Set(key string,val int)  {
	sm.rwmut.Lock()
	defer sm.rwmut.Unlock()
	sm.Map[key] = val
}
```

## WaitGroup

WaitGroup 常用在等goroutine结束。

```golang
package main
import (
	"fmt"
	"sync"
	"time"
)
func main()  {
  wg :=sync.WaitGroup{}
  wg.Add(2) //这边说明有2个任务
  go func() {
	defer wg.Done() //代表任务1结束
  	time.Sleep(1e9)
	fmt.Println("after 1 second")
  }()

  go func() {
	defer wg.Done()//代表任务2结束
	time.Sleep(2e9)
	fmt.Println("after 2 second")
  }()
  wg.Wait() //等待所有任务结束
  fmt.Println("the end")
}
```

```bash
after 1 second
after 2 second
the end
```


## 总结

本小节列举了goroutine使用传统的`Mutex`和`WaitGroup`做"线程"同步的例子，在go语言中，官方注重推荐使用`channel`做“线程”同步，下个小节我们着重介绍`channel`的使用。