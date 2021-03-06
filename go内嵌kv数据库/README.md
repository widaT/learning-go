# 为什么需要内嵌数据库

go的生态中有好多内嵌的k/v数据库,为什么我们需要内嵌数据库?

- 更高的效率:内嵌数据库因为和程序共享一个程序地址空间减少IPC开销，所以比独立数据库具有更高的性能
- 更简洁的部署方案:因为内嵌了,所以就不需要额外部署单独的数据库,减少程序依赖.
- 做单机存储引擎:一些优秀的内嵌数据库可以作为单机存储引擎,然后通过分布式一致性协议和分片策略可以集群成大型分布式数据库.例如etcd中使用boltDB,dgraph使用的badger.

说道内嵌型kv数据库,不得不提的就rocksdb了.这些年基于rocksdb单机存储引擎之上开发的分布式数据库数不胜数.例如tidb,cockroachdb等等.rocksdb是Facebook基于leveldb加强的kv存储引擎,用c++编写的,go需要通过cgo可以内嵌rocksdb,可以参考[gorocksdb](https://github.com/tecbot/gorocksdb)。但是go一般不用rocksdb作为内嵌kv数据库，首先rocksdb安装非常繁琐，二cgo的性能一直被诟病，而且容易造成内存泄露。


本小节依次介绍blotdb,goleveldb,badger三款表典型的内嵌数据库.

