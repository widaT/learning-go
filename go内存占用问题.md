# go程序内存占用问题
像golang和java这种带gc的语言，内存的申请，使用，和回收情况很多时候是不透明的。golang程序很容易看到rss（Resident Set Size 实际使用物理内存（包含共享库占用的内存））经常比实际用到的大很多的情况。这个时候不了解的原理很容易以为是内存泄露了。


先看一个程序
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

在系统负载比较低的时候，你会看到程序的Res 1G 左右，而且一直不变化。这个有点反直觉，我们向系统申请的50万个512维float32型的数组，后面实际上是已经删除了，
按理golang的gc应该回收这个1g的内存，然后归还给系统才对，可是你期望的事情没有发生。
幸好golang 只要设置环境变量 GODEBUG 就能看到gc的工作信息
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

我们可以在[go runtime](https://golang.org/pkg/runtime/)日志的格式信息
```
# gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
gc 10 @16.335s 0%: 0.004+1.0+0.004 ms clock, 0.016+0/1.0/2.9+0.016 ms cpu, 1006->1006->0 MB, 1784 MB goal, 4 P (forced)
```
- gc 10: gc的流水号，从1开始自增
- @16.335s: 从程序开始到当前打印是的耗时
- 0%: 程序开始到当前CPU时间百分比
- 0.004+1.0+0.004 ms clock: 0.004表示STW时间；1.0表示并发标记用的时间；0.004表示markTermination阶段的STW时间
- 0.016+0/1.0/2.9+0.016 ms cpu: for mark/scan are broken down in to assist time, background GC time, and idle GC time
- 1006->1006->0 MB: GC开始、GC结束、存活的heap大小
- 1784 MB goal:下次gc的目标值
- 4 P: 处理器数量
- (forced): 可能没有，代表程序中runtime.GC() 被调用


scvg: 1 MB released: gctrace的值大于0时，如果垃圾回收将内存返回给操作系统时，会打印一条summary,包括下一条数据

# 参考资料
[1](http://cbsheng.github.io/posts/godebug%E4%B9%8Bgctrace%E8%A7%A3%E6%9E%90/)
[2](https://golang.org/pkg/runtime/)
[3](https://colobu.com/2019/08/28/go-memory-leak-i-dont-think-so/)