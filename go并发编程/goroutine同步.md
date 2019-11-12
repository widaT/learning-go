# goroutine同步

## 锁（sync.Mutex）和读写锁（sync.RWMutex）

类似其他语言，golang也提供了锁和读写锁的的同步原语。

我们先看下go中锁（sync.Mutex）的使用。

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

我们在看下读写锁（sync.RWMutex），读写锁：读的时候不会阻塞读，会阻塞写；写的时候会阻塞写和读，我们可以用这个特性实现线程安全的map。

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

## channel
channel是goroutine同步绝对主角。

实现和WaitGroup一样的功能
```golang
func main() {
	ch := make(chan struct{}, 2)
	go func() {
		defer func() {
			ch <- struct{}{}
		}()
		time.Sleep(1e9)
		fmt.Println("after 1 second")
	}()

	go func() {
		defer func() {
			ch <- struct{}{}
		}()
		time.Sleep(2e9)
		fmt.Println("after 2 second")
	}()

	i := 0
	for _ = range ch {
		i++
		if i == 2 {
			close(ch)
		}
	}
	fmt.Println("the end")
}
```
这是个常用的1对n的生产消费模型，常常用于消费redis队列。
```golang
package main
import (
	"os"
	"fmt"
	"os/signal"
	"syscall"
	"sync"
)

var quit bool = false
const THREAD_NUM  = 5
func main() {
	sigs := make(chan os.Signal, 1)
	//signal.Notify 注册这个给定的通道用于接收特定信号。
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2)
	quitChan := make(chan bool)
	rowkeyChan := make(chan string, THREAD_NUM)
	go func() {
		 <-sigs  //等待信号
		quit = true
		close(rowkeyChan)
	}()
	var wg sync.WaitGroup
	for i := 0; i < THREAD_NUM;i++ {
		wg.Add(1)
		go func(n int) {
			for {
				rowkey,ok := <- rowkeyChan
				if !ok {
					break
				}
				//do something with rowkey
				fmt.Println(rowkey)
			}
			wg.Done()
		}(i)
	}
	go func() {
		wg.Wait()
		quitChan <- true
	}()

	for  !quit {
		//rowkey 可能来着redis的队列
		rowkey := ""
		rowkeyChan <- rowkey
	}
	<- quitChan
}

```

## 总结

本小节列举了goroutine常见的几个同步方式，goroutine不仅支持传统的“进程”同步方式，更重要的是通过channel的新型同步方式。下一小节我们将深入channel的运用。