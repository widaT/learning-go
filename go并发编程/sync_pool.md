# sync.Pool(临时对象池)

## 什么是sync.Pool

golang是带GC（垃圾回收）的语言，如果高频率的生成对象，然后有废弃这样子会给gc带来很大的负担，而且在go的内存申请上也会出现比较大的抖动。那么有什么办法减少gc负担，重用这些对象，然后又能让go的内存平缓些呢？答案是使用`sync.Pool`。

`sync.Pool`是golang用来存储临时对象的，这些对象通常是高频率生成销毁（这边还需要注意下，我们所说的对象是堆内存上的，而不是在栈内存上）。

## 使用sync.Pool

sync.Pool 有连个对外API

```
    func (p *Pool) Get() interface{} 
    func (p *Pool) Put(x interface{})
```
另外`sync.Pool`对象初始化的时候需要指定对象属性`New`是一个 `func() interface{}`方法类型，用来在没有可复用对象时重新生成对象。


其实`sync.Pool`的使用非常频繁，不管事go标准库还是第三方库都非常多的使用。在标准库`fmt`就使用到`sync.Pool`。

我们追踪下`fmt.Printf`的源码：

```go
func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}

func Sprintf(format string, a ...interface{}) string {
	p := newPrinter()
	p.doPrintf(format, a)
	s := string(p.buf)
	p.free()
	return s
}

var ppFree = sync.Pool{
	New: func() interface{} { return new(pp) },
}

func newPrinter() *pp {
	p := ppFree.Get().(*pp)
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
	ppFree.Put(p)
}
```