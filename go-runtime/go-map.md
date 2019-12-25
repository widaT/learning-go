# go map


go map的源码的目录为go/src/runtime/map.go


```go
type hmap struct {
	count     int // map的大小，len(map)返回的值
	flags     uint8 //状态flag 值有 1,2，4，8
	B         uint8  // bucket数量的对数值 例如B为5 则表明有2^5（32）个桶
	noverflow uint16 // 溢出桶（overflow buckets）的大概数量
	hash0     uint32 // hash seed
	buckets    unsafe.Pointer // 桶数组的指针，如果len为0，则指针为nil.
	oldbuckets unsafe.Pointer // 在扩容的时候这个指针不为nil
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)
	extra *mapextra // 当key和value都可内嵌的时候会用到这个字段
}
```

bucket的定义原始定义是：
```go
type bmap struct {
	tophash [bucketCnt]uint8
}
```
当时在实际编译后会变成如下结构，
```go
type bmap struct {
    tophash  [8]uint8
    keys     [8]keysize
    values   [8]keysize
    overflow uintptr
}
```













## map的几个特点

- go的map是非线程安全的，在并发读写的时候会触发`concurrent map read and map write`的panic。
- map的key不是有序的，所以在for range的时候经常看到map的key不是稳定排序的。
- map在删除某个key的时候不是真正的删除，只是标记为空，空间还在，所以在for range删除map元素是安全的。

 
## 参考资料

[map](https://github.com/cch123/golang-notes/blob/master/map.md)
[map的底层实现原理是什么.](https://github.com/qcrao/Go-Questions/blob/master/map/map%20%E7%9A%84%E5%BA%95%E5%B1%82%E5%AE%9E%E7%8E%B0%E5%8E%9F%E7%90%86%E6%98%AF%E4%BB%80%E4%B9%88.md)
