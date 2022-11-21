# 并发安全map —— sync.Map

go原生的map不是线程安全的，在并发读写的时候会触发`concurrent map read and map write`的panic。
`map`应该是属于go语言非常高频使用的数据结构。早期go标准库并没有提供线程安全的`map`，开发者只能自己实现，后面go吸取社区需求提供了线程安全的map——`sync.Map`。

`sync.Map` 提供5个如下api：
```go
    func (m *Map) Delete(key interface{})      //删除这个key的value
    func (m *Map) Load(key interface{}) (value interface{}, ok bool) //加载这个key的value
    func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) //原子操作加载，如果没有则存储
    func (m *Map) Range(f func(key, value interface{}) bool) //遍历kv
    func (m *Map) Store(key, value interface{}) //存储
```

## 使用sync.Map

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	sMap := sync.Map{}
	sMap.Store("a","b")
	ret,_:= sMap.Load("a")
	fmt.Printf("ret1 %t \n",ret.(string) == "b" )
	ret,loaded :=sMap.LoadOrStore("a","c")
	fmt.Printf("ret2 %t loaded:%t \n",ret.(string) == "b",loaded )
	ret,loaded =sMap.LoadOrStore("d","c")
	fmt.Printf("loaded %t \n",loaded)
	sMap.Store("e","f")
	sMap.Delete("e")
	sMap.Range(func(key, value interface{}) bool {
		fmt.Printf("k:%s v:%s \n", key.(string),value.(string))
		return true
	})
}
```

```bash
$ go run main.go
ret1 true 
ret2 true loaded:true 
loaded false 
k:a v:b 
k:d v:c 
```

## sync.Map底层实现

### `sync.Map`的结构体

```go
type Map struct {
	mu Mutex  //互斥锁保护dirty
	read atomic.Value //存读的数据，只读并发安全，存储的数据类型为readOnly
	dirty map[interface{}]*entry //包含最新写入的数据，等misses到阈值会上升为read
	misses int      //计数器，当读read的时候miss了就加+
}
```

```go
type readOnly struct {
    m  map[interface{}]*entry 
    amended bool //dirty的数据和这里的m中的数据有差异的时候true
}
```

从结构体上看Map有一定指针数据冗余，但是因为是指针数据，所以冗余的数据量不大。

### `Load`的源码：

```go
func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
	read, _ := m.read.Load().(readOnly) //读只读的数据
	e, ok := read.m[key]
	if !ok && read.amended { //如果没有读到且read的数据和drity数据不一致的时候
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly) //加锁后二次确认
		e, ok = read.m[key]
		if !ok && read.amended { //如果没有读到且 read的数据和drity数据不一致的时候
			e, ok = m.dirty[key]
			m.missLocked()     //misses +1,如果 misses 大于等于 m.dirty 则发送 read的值指向ditry
		}
		m.mu.Unlock()
	}
	if !ok {
		return nil, false
	}
	return e.load()
}
```

### `Store`的源码:

```go
func (m *Map) Store(key, value interface{}) {
	read, _ := m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok && e.tryStore(&value) { //如果在read中找到则尝试更新，tryStore中判断key是否已经被标识删除，如果已经被上传则更新不成功
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnly) //同上二次确认
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {// 如果entry被标记expunge，则表明dirty没有key，可添加入dirty，并更新entry。
			m.dirty[key] = e
		}
		e.storeLocked(&value)
	} else if e, ok := m.dirty[key]; ok { //如果dirty存在该key
		e.storeLocked(&value)
	} else { //key不存在
		if !read.amended { //read 和 dirty 一致
		  	// 将read中未删除的数据加入到dirty中
			m.dirtyLocked()
			 // amended标记为read与dirty不一致，因为即将加入新数据。
			m.read.Store(readOnly{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value)
	}
	m.mu.Unlock()
}
```

### `Delete`的源码：

```go
// Delete deletes the value for a key.
func (m *Map) Delete(key interface{}) {
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	if !ok && read.amended {   //在read中没有找到且 read和dirty不一致
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly) //加锁二次确认
		e, ok = read.m[key]
		if !ok && read.amended {
			delete(m.dirty, key)  //从dirty中删除
		}
		m.mu.Unlock()
	}
	if ok { //如果key在read中存在
		e.delete() //将指针置为nil，标记删除
	}
}
```
## 优缺点

优点: 通过read和dirty冗余的方式实现读写分离，减少锁频率来提高性能。
缺点：大量写的时候会导致read读不到数据而进一步加锁读取dirty，同时多次miss的情况下dirty也会频繁升级为read影响性能。 
因此`sync.Map`的使用场景应该是读多，写少。

## 总结

本小节介绍了`sync.Map`的使用，通过源码的方式了解`sync.Map`的底层实现，同时介绍了它的优缺点，以及使用场景。