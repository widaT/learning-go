# time

`time`包提供了时间的显示和测量用的函数。


## 获取当前时间

`time`中的`Now`方法会返回当前的时间（`time.Time`）。
```go
now := time.Now()
fmt.Println(now) //2020-10-23 11:02:53.356985487 +0800 CST m=+0.000042418

fmt.Println(now.Unix()) //获取时间戳  1603422283
fmt.Println(now.UnixNano())//获取纳秒时间戳  1603422283132651138


fmt.Println(now.Year()) //年
fmt.Println(now.Month()) //月
fmt.Println(now.Day()) //日

fmt.Println(now.Hour()) //时
fmt.Println(now.Minute()) //分
fmt.Println(now.Second()) //秒
fmt.Println(now.Nanosecond()) //纳秒
```

## 格式化输出

使用`time.Time`的`Format`方法格式化时间.
golang的时间格式比较有`特色`。`2006-01-02 15:04:05`代表`年-月-日 小时（24小时制）-分-秒`。


```go
fmt.Println(time.Now().Format("2006-01-02 15:04:05")) //2020-10-23 11:15:41
fmt.Println(time.Now().Format("2006-01-02 15:04:05.000")) //2020-10-23 11:15:41.439 带毫秒
fmt.Println(time.Now().Format("2006-01-02 15:04:05.000000")) //2020-10-23 11:15:41.439274 带微秒
fmt.Println(time.Now().Format("2006-01-02 15:04:05.000000000")) //2020-10-23 11:15:41.439277309 带纳秒
```

## 解析时间

时间戳转时间

```go
fmt.Println(time.Unix(1603422283, 0).Format("2006-01-02 15:04:05"))
```

字符串转时间

```go
t, err := time.Parse("2006-01-02 15:04:05", "2020-10-23 11:15:41")
fmt.Println(t)
```

使用`time.Date`构建时间

```go
t := time.Date(2020, 10, 23, 11, 15, 41, 0, time.Local)
fmt.Println(t)
```

## 时间测量

很多时候我们需要比较两个时间，甚至需要测量时间距离。


### 测耗时

`Time.Sub(t Time)`  方法测量自己和参数中的时间距离
`time.Slice(t Time)`函数测量参数时间t到现在的距离
`time.Until(t Time)`函数测量现在时间到参数t的距离

```go
start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
end := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)

difference := end.Sub(start)
fmt.Println(difference) //12h0m0s

fmt.Println(time.Since(start)) //182427h58m21.798590739s
fmt.Println(time.Until(end))   //-182415h58m21.798593974s
```

### Sleep函数

`Sleep`函数会让当前的`goroutine`休眠

```go
time.Sleep(d Duration)  //Duration  是int64的别名代办纳秒值 
```

```go
time.Sleep(1*time.Second) //sleep 1s
time.Sleep(1e9) //sleep 1s
```


## 使用`time.After`来处理超时

```go
func After(d Duration) <-chan Time
```

```go
package main

import (
	"fmt"
	"time"
)

var c chan int

func handle(int) {}

func main() {
	select {
	case m := <-c:
		handle(m) 
	case <-time.After(10 * time.Second):
		fmt.Println("timed out")
	}
}
```

## 定时器Time.Tick

使用`Time.Tick`会产生一个定时器.

```go
func Tick(d Duration) <-chan Time
```

```go
package main

import (
	"fmt"
	"time"
)

func statusUpdate() string { return "" }

func main() {
	c := time.Tick(5 * time.Second) //每5s 生产时间 往channel 发送
	for next := range c {
		fmt.Printf("%v %s\n", next, statusUpdate())
	}
}
```