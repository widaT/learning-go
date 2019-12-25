# go runtime简介

有别于java和c#的runtime是一个虚拟机,go的runtime和我们的go代码一同编译成二进制可执行文件。

go的runtime负责：

- goroutine的调度
- go内存分配
- 垃圾回收（GC）
- 封装了操作系统底层操作，如syscall，原子操作，CGO
- map，channel，slice,string内置类型的实现
- 反射（reflection）的实现
- pprof，trace，race的实现



## 参考资料
 

[深度长文：深入Go Runtime，你能做到浅出吗?](https://mp.weixin.qq.com/s/ivO-USpxiyrL-9BzgE8Vcg)