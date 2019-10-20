# go程序物理内存占用高的问题
最近手头上的一个go项目，go版本是go1.13.1。项目上线后发现rss（Resident Set Size实际使用物理内存）很快就飚7g（服务器内存8g），而且很长时间内存都不下降。一开始还以为的内存泄露了。后来经过一番折腾才发现这个不是内存泄露。

## GODEBUG
先写一个模拟程序
```go
package  main
import (
	"runtime/debug"
	"time"
)
func main()  {
	num := 500000
  	var bigmap = make(map[int]*[512]float32)
  	for i := 0;i < num;i++ {
  		bigmap[i] = &[512]float32{float32(i)}
	}

  	println(len(bigmap))
  	time.Sleep(15e9)
	for i := 0;i < num;i++ {
		delete(bigmap,i)
	}

	debug.FreeOSMemory()
	time.Sleep(1000e9)
}
```

```bash
# go run main.go
```

然后打开终端
```bash
# top -p pid
```
![1](../img/1.png)

在系统负载比较低的时候，你会看到程序的Res 1G 左右，而且一直不变化。这个有点反直觉，我们向系统申请的50万个512维float32型的数组，后面实际上是已经删除了，
按理说golang的gc应该回收这个1g的内存，然后归还给系统才对，可是这样子的事情并没有发生。
到底发生了什么？我们去追踪下gc日志。
```bash
# GODEBUG=gctrace=1 go run main.go
...
500000
gc 10 @16.335s 0%: 0.004+1.0+0.004 ms clock, 0.016+0/1.0/2.9+0.016 ms cpu, 1006->1006->0 MB, 1784 MB goal, 4 P (forced)
scvg: 1 MB released
scvg: inuse: 835, idle: 188, sys: 1023, released: 17, consumed: 1005 (MB)
forced scvg: 1005 MB released
forced scvg: inuse: 0, idle: 1023, sys: 1023, released: 1023, consumed: 0 (MB)
```
我们可以在[go runtime](https://golang.org/pkg/runtime/) 看到gc日志的格式的相关信息

## gc日志
格式如下：
` gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P`

```
gc 10 @16.335s 0%: 0.004+1.0+0.004 ms clock, 0.016+0/1.0/2.9+0.016 ms cpu, 1006->1006->0 MB, 1784 MB goal, 4 P (forced)
```
- gc 10: gc的流水号，从1开始自增
- @16.335s: 从程序开始到当前打印是的耗时
- 0%: 程序开始到当前CPU时间百分比
- 0.004+1.0+0.004 ms clock: 0.004表示STW时间；1.0表示并发标记用的时间；0.004表示markTermination阶段的STW时间
- 0.016+0/1.0/2.9+0.016 ms cpu: 0.016表示整个进程在mark阶段STW停顿时间;0/1.0/2.9有三块信息，0是mutator assists占用的时间，2.9是dedicated mark workers+fractional mark worker占用的时间，2.9+是idle mark workers占用的时间。0.016 ms表示整个进程在markTermination阶段STW停顿时间(0.050 * 8)。
- 1006->1006->0 MB: GC开始、GC结束、存活的heap大小
- 1784 MB goal:下次gc的目标值
- 4 P: 处理器数量
- (forced): 可能没有，代表程序中runtime.GC() 被调用

## scvg日志: 
go语言把内存归还给操作系统的过程叫scavenging，scvg日志记录的就是这个过程的日志信息。
scvg 每次会打印两条日志
格式：
```
scvg#: # MB released  printed only if non-zero
scvg#: inuse: # idle: # sys: # released: # consumed: # (MB)
```
```
scvg: 1 MB released
scvg: inuse: 835, idle: 188, sys: 1023, released: 17, consumed: 1005 (MB)
```
- 1 MB released: 返回给系统1m 内存
- inuse: 正在使用的
- idle ：空闲的
- sys ： 系统映射内存
- released:已经归还给系统的内存
- consumed：已经向操作系统申请的内存

所以从上面的介绍来看，最后的归还日志
```
forced scvg: inuse: 0, idle: 1023, sys: 1023, released: 1023, consumed: 0 (MB)
```
说明go语言已经正常把内存交给操作系统了了。可是RSS信息显示go仍然占用这这个内存。
这又是为什么？
在[go runtime](https://golang.org/pkg/runtime/)中我们找到一段如下的文字。
```
madvdontneed: setting madvdontneed=1 will use MADV_DONTNEED
instead of MADV_FREE on Linux when returning memory to the
kernel. This is less efficient, but causes RSS numbers to drop
more quickly.
```
翻译一下：

```
madvdontneed：如果设置GODEBUG=madvdontneed=1，golang归还内存给操作系统的方式将使用MADV_DONTNEED，而不是Linux上的MADV_FREE的方式。虽然MADV_DONTNEED效率较低，但会程序RSS下降得更快速。
```

动手试一试
```bash
# GODEBUG=madvdontneed=1 go run main.go
```
![2](../img/2.png)
这下RSS正常了。


## MADV_FREE和MADV_DONTNEED

MADV 是linux kernel的内存回收的方式，有两种变体（MADV_FREE和MADV_DONTNEED，可以参考
[MADV_FREE functionality](http://lkml.iu.edu/hypermail/linux/kernel/0704.3/3962.html)文档。

### 简单来说
- MADV_FREE将内存页（page）延迟回收。当内核内存紧张时,这些内存页将会被优先回收，如果应用程序在页回收后又再次访问，内核将会返回一个新的并设置为0的页。而如果内核内存充裕时，标识为MADV_FREE的页会仍然存在，后续的访问会清掉延迟释放的标志位并正常读取原来的数据。
- MADV_DONTNEED告诉内核这段内存今后"很可能"用不到了，其映射的物理内存尽管拿去好了，因此物理内存可以回收以做它用，但虚拟空间还留着，下次访问时会产生缺页中断，这个缺页中断会触发重新申请物理内存的操作。

go1.12之后采用了MADV_FREE的方式，这样子个go程序和内核间的内存交互会更加高效，但是带来的一个副作用就是rss下降的比较慢（系统资源充裕时，甚至你感觉不到他是在下降）。在关注 `top` 信息的时候，别只关注RES和%MEM，还要关注下buff/cache。在系统资源紧张的时候MADV_FREE方式标记的内存也会很快回收，并不需要太担心。另外使用`GODEBUG=madvdontneed=1`会强制使用原先MADV_DONTNEED的方式。



## 补充
- golang的内存申请方式会比较激进，像map扩容的时候一次申请的比实际需要的多得多。然后在delete map的时候内存不会释放。
- 如果在程序中需要反复申请对象，然后销毁的话，应该使用`sync.pool`来重复使用申请过的内存，这样子可以让程序申请的系统内存相对平稳。

# 参考资料
- [GODEBUG之gctrace解析](http://cbsheng.github.io/posts/godebug%E4%B9%8Bgctrace%E8%A7%A3%E6%9E%90/)
- [go runtime](https://golang.org/pkg/runtime/)
- [Go内存泄漏？不是那么简单!](https://colobu.com/2019/08/28/go-memory-leak-i-dont-think-so/)
- [MADV_FREE functionality](http://lkml.iu.edu/hypermail/linux/kernel/0704.3/3962.html)