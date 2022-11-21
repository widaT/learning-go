# sync.Pool(临时对象池)

## 什么是sync.Pool

golang是带GC（垃圾回收）的语言，如果高频率的生成对象，然后有废弃这样子会给gc带来很大的负担，而且在go的内存申请上也会出现比较大的抖动。那么有什么办法减少gc负担，重用这些对象，然后又能让go的内存平缓些呢？答案是使用`sync.Pool`。

`sync.Pool`是golang用来存储临时对象的，这些对象通常是高频率生成销毁（这边还需要注意下，我们所说的对象是堆内存上的，而不是在栈内存上）。

## 使用sync.Pool

sync.Pool 有两个对外API

```
    func (p *Pool) Get() interface{} 
    func (p *Pool) Put(x interface{})
```
另外`sync.Pool`对象初始化的时候需要指定属性`New`是一个 `func() interface{}`函数类型，用来在没有可复用对象时重新生成对象。

其实`sync.Pool`的使用非常频繁，不管事go标准库还是第三方库都非常多的使用。
在标准库`fmt`就使用到`sync.Pool`，我们追踪下`fmt.Printf`的源码：

```go
func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}

func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	p := newPrinter()
	p.doPrintf(format, a)
	n, err = w.Write(p.buf)
	p.free()zongj
	return
}

var ppFree = sync.Pool{
	New: func() interface{} { return new(pp) }, //指定生成对象函数
}

func newPrinter() *pp {
	p := ppFree.Get().(*pp)    //从pool中获取可复用对象，如果没有对象池会重新生成一个，注意这边拿到对象后会reset对象
	p.panicking = false
	p.erroring = false
	p.wrapErrs = false
	p.fmt.init(&p.buf)
	return p
}

func (p *pp) free() {
	if cap(p.buf) > 64<<10 {
		return
	}
	p.buf = p.buf[:0]
	p.arg = nil
	p.value = reflect.Value{}
	p.wrappedErr = nil
	ppFree.Put(p)             //用完后重新放到pool中
}
```

从上面的案例中大概可以看出`sync.Pool`是如何使用的。接下来我们写一个demo程序，看下另外一个`sync.Pool`的高频使用场景


```go
package main

import (
	"io"
	"log"
	"net"
	"sync"
)

func main() {
	bufpool := sync.Pool{}
	bufpool.New = func() interface{} {
		return make([]byte, 32768)
	}
	Pipe := func(c1, c2 io.ReadWriteCloser) {
		b := bufpool.Get().([]byte)
		b2 := bufpool.Get().([]byte)
		defer func() {
			bufpool.Put(b)
			bufpool.Put(b2)
			c1.Close()
			c2.Close()
		}()

		go io.CopyBuffer(c1, c2, b)
		io.CopyBuffer(c2, c1, b2)
	}
	l,err := net.Listen("tcp",":9999")
	if err !=nil {
		log.Fatal(err)
	}
	for  {
		conn,err := l.Accept()
		if err !=nil {
			log.Fatal(err)
		}
		client ,err:= net.Dial("tcp","127.0.0.1:80")
		if err !=nil {
			log.Fatal(err)
		}
		go Pipe(conn,client)
	}
}
```
这是个的代理程序，任何连到本机9999端口的tcp链接都会转发到本地的80端口。我们使用`io.CopyBuffer`实现数据双工互相拷贝。
`io.CopyBuffer`会频繁使用到缓存`[]byte`对象，我们用`sync.Pool`重复使用`[]byte`.

运行一下程序
```bash
$ go run main.go
$ curl http://localhost:9999
```

## 总结

本小节介绍了`sync.Pool`的使用方式。`sync.Pool`能减轻go GC的负担，同时减少内存的分配，是保障go程序内存分配平缓的重要手段。