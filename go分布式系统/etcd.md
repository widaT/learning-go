# Etcd

在本章节我们不止一次提到etcd，因为etcd这个项目在go语言分布式生态里头非常的重要。etcd个高可用的分布式键值(key-value)数据库，采用raft协议作为一致性算法，基于Go语言实现。etcd在go分布式系统生态里头的位置相当于zookeeper在java分布式系统生态的位置。tidb的raft系统是基于etcd的实现然后用rust实现的，tidb的pd系统还直接内嵌了etcd。在微服务生态里头etcd可以做服务发现组件使用，kubernetes也用etcd保存集群所有的网络配置和对象的状态信息。


## etcd的特点

- 简单：提供定义明确的面向用户的API（有http和grcp实现，注意http和grpc分别用了两套存储，数据不同）
- 安全：支持client端证书的TLS
- 快速：基准测试支持10,000次write/s
- 可靠：使用raft协议保证数据一致性

对于想研究分布式系统的gopher来说，etcd是一个非常值得好好阅读源码的项目。它的raft实现，mvcc，wal实现都非常典型。

## etcd使用

etcd项目有两个go client，一个基于http restful对于源码目录里头的client，另外一个基于grcp对于源码目录的clientv3，需要注意是这两个api操作的数据是不一样的。client支持的功能较小，而且性能也比较有问题，在go生态里头都会用基于grpc的clientv3。