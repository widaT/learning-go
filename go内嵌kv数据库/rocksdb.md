# Rocksdb

其实rocksdb在go生态中使用的不是很多，但是说道内嵌型kv数据库，rocksdb绝对是目前最流行的。所以本小节我们介绍下rocksdb。
由Facebook基于levelDB开发，同时针对Flash存储进行优化，延迟极小，同样支持ACID，rocksdb还借鉴了hbase的很多理念，比如支持column families。


## LSM-Tree

LSM Tree（Log-Structured Merge-Tree），是为了解决在内存不足，磁盘随机IO太慢下的写入性能问题。在LSM中增删改操作都是新增记录，新的记录会最早被索引到，这样才的操作会造成数据重复，后续的合并操作会消除这样的冗余。LSM通过这样的方式把随机IO变成顺序IO，大大提高写入性能。

了解LSM细节可以查看[LSM树原理探究](https://juejin.im/post/6844903863758094343)。

