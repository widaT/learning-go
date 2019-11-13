# 使用channel做goroutine同步

channel和goroutine是golang CSP并发模型的承载体，是goroutine同步绝对主角。channel的使用场景远比`Mutex`和`WaitGroup`多得多。

## 实现和WaitGroup一样的功能

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

## 生产消费模型

这是个常用的1对N的生产消费模型，常常用于消费redis队列。

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

上面的代码稍微修改下很容易支持N：M的生产消费模型，这边就不再赘述。

## 赛道模型

在现实很多场景我们需要并发的做个任务，我们想知道他们的先后顺序，或者只想知道最快的那个。channel的特性很容易做到这个需求
```golang
package main
import (
	"fmt"
	"net/http"
)
func main() {
	ch := make(chan string)
	go func() {
		http.Head("https://www.taobao.com")
		ch <- "taobao"
	}()

	go func() {
		http.Head("https://www.jd.com")
		ch <- "jd"
	}()
	go func() {
		http.Head("https://www.vip.com")
		ch <- "vip"
	}()

	//只想知道最快的
	dm := <- ch
	fmt.Println(dm)

/* 知道它们的排名
	for dm:= range ch {
		fmt.Println(dm)
	}
	*/
}
```

## 如何控制并发数

虽然golang新开一个goroutine占用的资源很小，但是无谓的goroutine开销对golang的调度性能很有影响，同时也会浪费cpu的资源。那么golang中如何控制并发数。

### 方式1——使用带缓冲的channel

```golang
package main
import (
	"fmt"
	"time"
)

type Task struct {
	Id int
}

func work(task *Task,limit chan struct{})  {
	time.Sleep(1e9)
	fmt.Println(task.Id)
	<- limit
}
func main() {
	ch := make(chan *Task,10)
	limit := make(chan struct{},3)
	go func() {
		for i:=0;i<10;i++ {
			ch <- &Task{i}
		}
		close(ch)
	}()

	for {
		task ,ok :=<- ch
		if ok {
			limit <- struct{}{}
			go work(task,limit)
		}
	}
}
```

### 方式2——使用master-worker模型，worker数量固定

```golang
package main

import (
	"fmt"
	"time"
)
const WorkerNum = 3
type Task struct {
	Id int
}
func work(ch chan *Task)  {
	for {//worker中循环消费任务
		task,ok :=<- ch          
		if ok {
			time.Sleep(1e9)
			fmt.Println(task.Id)
		}
	}
}
func main() {
	ch := make(chan *Task,3)
	exit := make(chan struct{})
	go func() {
		for i:=0;i<10;i++ {
			ch <- &Task{i}
		}
		close(ch)
	}()

	for i:=0;i<WorkerNum;i++{  //控制worker数量
		go work(ch)
	}
	<- exit
}
```


## 控制goroutine优雅退出

除了我们需要控制goroutine的数量之外，我们还需要控制goroutine的生命周期同样以防止不不要的资源和性能消耗。那么接下来我们介绍下控制goroutine生命周期的办法。

### 方法1——goroutine超时控制

```golang
package main
import (
	"fmt"
	"time"
)
func main() {
	ch := make(chan struct{})
	task := func() {
		time.Sleep(3e9)
		ch <- struct{}{}
	}
	go task()
	select {
		case <-ch: //在3s内正常完成
			fmt.Println("task finish")
		case <-time.After(3*time.Second): //超过3秒
			fmt.Println("timeout")
	}
}
```

```bash
$ go run main.go 
timeout
```

我们使用了 `slecet` 和 `time.After` 来控制goroutine超时。

### 方法2——使用`context.Context`

```golang
package main
import (
	"context"
	"fmt"
	"time"
)
func main() {
	ctx ,cancel:= context.WithCancel(context.Background())  //使用带cancel函数的context
	task := func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done(): //cancel函数已经被执行了
			default:
				time.Sleep(1e9)
				fmt.Println("hello")
			}
		}
	}
	go task(ctx)
	time.Sleep(3e9) //等待3s
	cancel()          //让goroutine退出
}
```

```bash
$ go run main.go
hello
hello
```

当然我们还可以使用 `context.WithTimeout`的方式

```golang
package main

import (
	"context"
	"fmt"
	"time"
)
func main() {
	exit := make(chan struct{})
	ctx,_ := context.WithTimeout(context.Background(),3*time.Second) //使用带超时的context
	task := func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done(): //cancel函数已经被执行了
			default:
				time.Sleep(1e9)
				fmt.Println("hello")
			}
		}
		exit<- struct{}{}
	}
	go task(ctx)
	<-exit
}
```

```bash
$ go run main.go 
hello
hello
hello
```

这样的用法和上面的用法的区别是可以不需要手动去调用cancel函数。

## 总结

本小节我们介绍了使用`channel`同步goroutine的方法，介绍了`channel`的常见使用场景，已经如果控制`goroutine`的几种方法。这些方法在很是实际开发中都会用到。