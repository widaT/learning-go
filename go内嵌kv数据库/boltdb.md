# [BoltDB](https://github.com/etcd-io/bbolt)

我们现在说的bolt一般都会说是etcd维护的bolt版本,老的bolt作者已经不再维护。bolt提供轻量级的kv数据库，支持完整的ACID事务。
bolt是基于内存映射（mmap）的数据库，通常情况使用的内存会比较大，bolt适合kv数量不是特别大的场景。

# 创建数据库

```go
package main

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

func main() {
	// 打开my.db 如果文件不存在则会创建
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}
```

# 事务

bolt支持只读事务和读写事务， 在同一个时间只能有一个读写事务，但是可以有多个只读事务。

## 只读事务使用 `db.View`

```go
err := db.View(func(tx *bolt.Tx) error {
	...
	return nil
})
```

## 读写事务使用 `db.Update`

```go

err := db.Update(func(tx *bolt.Tx) error {
	...
	return nil
})
```

## Buckets(桶)

Bolt在内部用Buckets来组织相关的kv键值对，在同一个Bucket里头key值是不允许重复的。


### 操作Bucket

```go
db.Update(func(tx *bolt.Tx) error {
	b, err := tx.CreateBucket([]byte("MyBucket")) //创建bucket，通常会使用 CreateBucketIfNotExists
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
    }
    
    err := b.Put([]byte("answer"), []byte("42")) //写入kv
    v := b.Get([]byte("answer")) //获取value
    fmt.Printf("%s",v)
	return nil
})
```

## key的遍历

在kv数据库中我们需要对key进行精心设计，无论在取值或者遍历的时候都需要快速的定位key的位置。在bolt中key是基于B树有序的。
一般有如下三种场景遍历key，遍历，范围遍历，前缀遍历

### 遍历桶中所有key

```go
db.View(func(tx *bolt.Tx) error {
	b := tx.Bucket([]byte("MyBucket"))
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		fmt.Printf("key=%s, value=%s\n", k, v)
	}
	return nil
})
```

或者使用 `ForEach` 遍历所有桶下面所有key

```go
db.View(func(tx *bolt.Tx) error {
	// Assume bucket exists and has keys
	b := tx.Bucket([]byte("MyBucket"))

	b.ForEach(func(k, v []byte) error {
		fmt.Printf("key=%s, value=%s\n", k, v)
		return nil
	})
	return nil
})
```

### key范围遍历

```go
db.View(func(tx *bolt.Tx) error {
	c := tx.Bucket([]byte("Events")).Cursor()

	min := []byte("1990-01-01T00:00:00Z")
	max := []byte("2000-01-01T00:00:00Z")

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		fmt.Printf("%s: %s\n", k, v)
	}

	return nil
})
```

### key前缀遍历

```go
db.View(func(tx *bolt.Tx) error {
	c := tx.Bucket([]byte("MyBucket")).Cursor()
	prefix := []byte("1234")
	for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
		fmt.Printf("key=%s, value=%s\n", k, v)
	}

	return nil
})
```