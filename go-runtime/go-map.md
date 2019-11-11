# go map

## map的几个特点
- go的map是非线程安全的，在并发读写的时候会触发`concurrent map read and map write`的panic。
- map的key不是有序的，所以在for range的时候经常看到map的key不是稳定排序的。
- map在删除某个key的时候不是真正的删除，只是标记为空，空间还在，所以在for range删除map元素是安全的。

 