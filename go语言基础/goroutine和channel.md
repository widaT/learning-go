# goroutine和channel


## goroutine
golang原生支持并发，在golang中每一个并发单元叫`goroutine`。`goroutine`你可以理解为golang实现轻量级的用户态线程。go程序启动的时候
其主函数就开始在一个单独的goroutine中运行，这个goroutine我们叫main goroutine。

在golang中启动一个goroutine的成本很低，通常只需要在普通函数执行前加`go`关键字就可以。

```go
func fn(){
    fmt.Println("hello world")
} 

fn()
go fn() //启用新的goroutine中执行 fn函数
```


## channel

channel是golang中goroutine通讯的一直机制，每一个channel都带有类型，channel使用`make`来创建。

```go
ch1 :=make(chan int) //创建不带缓冲区的channel
ch2 :=make(chan int，10)//创建带缓冲区的channel
```

channel是指针类型，它的初始值为`nil`，channel的操作主要有 发送，接受，和关闭。

```go
ch <- x  // 发送x
x = <-ch // 从channel中接收 然后赋值给x
<-ch     // 接收后丢弃
close(x) //
```

不带缓存的`channel`，当队列有一个值的时候，再往`channel`发送数据时会阻塞这个`goroutine`直到`channel`可以写入。
带缓存的`channel`，当这个channel的缓存队列没有满是往channel写数据是不会阻塞的，当队列满是会阻塞这个`goroutine`。
`channel`读取的时候如果`channel`队列有值会读取，队列为空的时候会塞这个`goroutine`直到`channel`有值可以读取。
当一个`channel`被close后，基于该channel的发送操作都将导致`panic`，接收操作可以接受到已经channel队列里头数据，channel队列为空时产生一个零值的数据。


golang中的`channel`还可以带方向。

```go
var out chan<- int  //只发送，不能接收
var in <-chan int   //只接收，不能发送 注意 对只接收的channel close会引起编译错误
```


本小节简介了goroutine和channel相关概念，具体的并发编程模型将在下个章节详细探讨。