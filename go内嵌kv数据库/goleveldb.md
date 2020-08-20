# Leveldb

LevelDB是Google开源的由C++编写的基于LSM-Tree的KV数据库，有超高的随机写，顺序读/写性能，但是随机读的性能很一般，LevelDB一般应用在查询少写多的场景。


## LSM-Tree

LSM Tree（Log-Structured Merge-Tree），是为了解决在内存不足，磁盘随机IO太慢下的写入性能问题。在LSM中增删改操作都是新增记录，新的记录会最早被索引到，这样才的操作会造成数据重复，后续的合并操作会消除这样的冗余。LSM通过这样的方式把随机IO变成顺序IO，大大提高写入性能。

了解LSM-Tree细节可以查看[LSM树原理探究](https://juejin.im/post/6844903863758094343)。


## goleveldb使用


### 创建db

```go
db, err := leveldb.OpenFile("path/to/db", nil)
defer db.Close()
```

### 增删改

```go
data, err := db.Get([]byte("key"), nil)
...
err = db.Put([]byte("key"), []byte("value"), nil)
...
err = db.Delete([]byte("key"), nil)
```


### 迭代器


#### 遍历所有

```go
iter := db.NewIterator(nil, nil)
for iter.Next() {
	key := iter.Key()
	value := iter.Value()
	...
}
iter.Release()
err = iter.Error()
```

#### 设置起始点遍历

```go
iter := db.NewIterator(nil, nil)
for ok := iter.Seek(key); ok; ok = iter.Next() {
	// Use key/value.
	...
}
iter.Release()
err = iter.Error()
```

#### 范围遍历

```go
iter := db.NewIterator(&util.Range{Start: []byte("foo"), Limit: []byte("xoo")}, nil)
for iter.Next() {
	// Use key/value.
	...
}
iter.Release()
err = iter.Error()
...
```

#### 前缀遍历

```go
iter := db.NewIterator(util.BytesPrefix([]byte("foo-")), nil)
for iter.Next() {
	// Use key/value.
	...
}
iter.Release()
err = iter.Error()
```

### 批量操作

```go
batch := new(leveldb.Batch)
batch.Put([]byte("foo"), []byte("value"))
batch.Put([]byte("bar"), []byte("another value"))
batch.Delete([]byte("baz"))
err = db.Write(batch, nil)
```


### 使用布隆过滤器

```
布隆过滤器使用极少空间代价（bitmap）来判断一个值是否在这个集合里，它使用多个hash函数对key计算然后映射到bitmap上的位置，如果这个位置的值为1，则代表这个值可能存在（注意是可能存在），如果这个位置值为0，代表这个key一定不存在这个集合里。如果哈希函数够多，我们判断是否存在的可信度就越高。当然某些情况下判断出来值存在的误判，在实际应用场景中我们无非耗费点资源去实际的kv数据库里头查下看否能拿到value。布隆过滤器能快速的知道key不存这样子可以减少大量的好资源操作。
```

```go
o := &opt.Options{
	Filter: filter.NewBloomFilter(10),
}
db, err := leveldb.OpenFile("path/to/db", o)
...
defer db.Close()
...
```