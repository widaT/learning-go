# [Badger](https://github.com/dgraph-io/badger)

badger 同样是基于LSM tree,但是不同于levevdb他的key/value在存储的时候是分离的,特别时候value比较大的场景.


## badger使用

### 创建db

```go
package main
import (
	"log"
	badger "github.com/dgraph-io/badger/v2"
)

func main() {
  db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
  if err != nil {
	  log.Fatal(err)
  }
  defer db.Close()
}
```

### 增删改

```go

//写数据
err := db.Update(func(txn *badger.Txn) error {
  err := txn.Set([]byte("answer"), []byte("42"))
  return err
})



//读数据
err := db.View(func(txn *badger.Txn) error {
  item, err := txn.Get([]byte("answer"))
  //handle(err)
  var valNot []byte
  var valCopy []byte
  //因为key-value是分类 bager获取value的方式比较独特
  err := item.Value(func(val []byte) error {
    fmt.Printf("The answer is: %s\n", val)
    // 拷贝val 或者解析val是可以
    valCopy = append([]byte{}, val...)
    valNot = val //不要这么干
    return nil
  })
  //handle(err)

  // 不要这么干，这样会经常导致bug
  fmt.Printf("NEVER do this. %s\n", valNot)

  // You must copy it to use it outside item.Value(...).
  fmt.Printf("The answer is: %s\n", valCopy)

  // 或者你也可以这么干
  valCopy, err = item.ValueCopy(nil)
  //handle(err)
  fmt.Printf("The answer is: %s\n", valCopy)
  return nil
})

//删除数据
err := db.Update(func(txn *badger.Txn) error {
  err := txn.Delete([]byte("answer"))
  return err
})

```

### 事务
 
入上面的示例代码，badger的使用有只读事务和可写事务分别用

```go
//只读事务
err := db.View(func(txn *badger.Txn) error {
  return nil
})
```

```go
//可写事务
err := db.Update(func(txn *badger.Txn) error {
  return nil
})
```

一般情况下`db.View`和`db.Update`大多情况已经满足大多数场景.但是默写情况下你可能需要自己管理事务的开启和关闭.badger提供了`DB.NewTransaction()`,`Txn.Commit()`,`Txn.Discard()`三个方法来手动管理事务.


```go
// 开启可写事务
txn := db.NewTransaction(true)

//别忘了用txn.Discard()清理事务
defer txn.Discard()


err := txn.Set([]byte("answer"), []byte("42"))
if err != nil {
    return err
}

// 提交事务或者出来事务
if err := txn.Commit(); err != nil {
    return err
}
```

还需要注意的是在badger事务大小是受限制的，我们在处理事务的是需要出来事务的`error`，特别需要关注`badger.ErrTxnTooBig`，我需要分批处理下，处理方式入下代码。

```go
updates := make(map[string]string)
txn := db.NewTransaction(true)
for k,v := range updates {
  if err := txn.Set([]byte(k),[]byte(v)); err == badger.ErrTxnTooBig {
    _ = txn.Commit()
    txn = db.NewTransaction(true)
    _ = txn.Set([]byte(k),[]byte(v))
  }
}
_ = txn.Commit()
```

### 遍历

只遍历key

```go
err := db.View(func(txn *badger.Txn) error {
  opts := badger.DefaultIteratorOptions
  opts.PrefetchValues = false
  it := txn.NewIterator(opts)
  defer it.Close()
  for it.Rewind(); it.Valid(); it.Next() {
    item := it.Item()
    k := item.Key()
    fmt.Printf("key=%s\n", k)
  }
  return nil
})
```

遍历key-value

```go
err := db.View(func(txn *badger.Txn) error {
  opts := badger.DefaultIteratorOptions
  opts.PrefetchSize = 10
  it := txn.NewIterator(opts)
  defer it.Close()
  for it.Rewind(); it.Valid(); it.Next() {
    item := it.Item()
    k := item.Key()
    err := item.Value(func(v []byte) error {
      fmt.Printf("key=%s, value=%s\n", k, v)
      return nil
    })
    if err != nil {
      return err
    }
  }
  return nil
})
```

前缀遍历

```go
err := db.View(func(txn *badger.Txn) error {
  it := txn.NewIterator(badger.DefaultIteratorOptions)
  defer it.Close()
  prefix := []byte("1234")
  for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
    item := it.Item()
    k := item.Key()
    err := item.Value(func(v []byte) error {
      fmt.Printf("key=%s, value=%s\n", k, v)
      return nil
    })
    if err != nil {
      return err
    }
  }
  return nil
})
```


### value文件清理

由于`badger`的`key`和`value`是分开存储的，`key`存储在`LSM-TREE`中，而value则独立于`LSM-TREE`之外，`LSM-TREE`的压缩过程，不会涉及value的合并。另外badger的事务基于MVCC的所以value是存在很多个版本的。总的来说手动清理value文件是必须的。

badger提供`db.RunValueLogGC`来清理value文件。
通常我们需要单独启用一个goroutine来执行。
例子如下：
```go
ticker := time.NewTicker(5 * time.Minute)
defer ticker.Stop()
for range ticker.C {
again:
    err := db.RunValueLogGC(0.7)
    if err == nil {
        goto again
    }
}
```