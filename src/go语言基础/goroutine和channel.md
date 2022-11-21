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
x := <-ch // 从channel中接收 然后赋值给x
<-ch     // 接收后丢弃
close(x) //
```

不带缓存的`channel`，发送的时候会阻塞当前`goroutine`，知道`channel`的信息被其他`goroutine`消费。
带缓存的`channel`，当这个channel的缓存队列没有满是往`channel`写数据是不会阻塞的，当队列满是会阻塞这个`goroutine`。
`channel`读取的时候如果`channel`队列有值会读取，队列为空的时候会塞这个`goroutine`直到`channel`有值可以读取。
当一个`channel`被close后，基于该channel的发送操作都将导致`panic`，接收操作可以接受到已经channel队列里头数据，channel队列为空时产生一个零值的数据。


golang中的`channel`还可以带方向。

```go
var out chan<- int  //只发送，不能接收
var in <-chan int   //只接收，不能发送 注意 对只接收的channel close会引起编译错误
```


## select控制结构和channel

`select`是golang中的一个控制结构，和switch表达式有点相似，不同的是`select`的每个`case`分支必须是一个通信操作（发送或者接收）。
`select`随机选择可执行`case`分支。如果没有`case`分支可执行，它将阻塞，直到有`case`分支可执行。
带`default`分支的`select`在没有可执行`case`分支时会执行`default`。

```go
ch :=make(chan int,1)
select {
    case ch <- 1  :
     //代码     
    case n:= <-ch  :
       //代码   
    default : //default 是可选的
       //代码   
}
```

`select`语句的特性：
- 每个`case`分支都是通信操作
- 所有`case分支`表达式都会被求值
- 如果任意某个`case分支`不阻塞，它就执行，其他被忽略。
- 如果有多个`case`分支都不阻塞，`select`会随机地选出一个执行。其他不会执行。
    否则：
    1. 如果有 `default` ，则执行该语句。
    2. 如果没有 `default`，`select` 将阻塞，直到某个`case分支`可以运行；


```go
ch := make(chan int, 1) //这边需要使用1个缓冲区，这样子可以在一个goroutine内使用

for {
	select {
	case ch <- 1:
		fmt.Println("send")
	case n := <-ch:
		fmt.Println(n)
	default:
		fmt.Println("dd")
	}
	time.Sleep(1e9)
}
```

本小节简介了`goroutine`和`channel`相关概念，具体的并发编程模型将在下个章节详细探讨。