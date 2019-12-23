# 分布式一致性算法——Gossip算法

前两节我们介绍了CA算法的典型——2PC，CP算法的典型——Raft，本小节要介绍AP算法的典型——Gossip算法。

Gossip算法有很多别名，流言算法，流行病算法等，Gossip最早是在 1987 年发表在 ACM 上的论文 《Epidemic Algorithms for Replicated Database Maintenance》中被提出。后来因为Cassandra中运用名声大噪。近几年区块链盛行，比特币（Bitcoin）和超级记账本（ Fabric）同样使用了Gossip。

我们看下Gossip的算法模型：

和Raft算法不同，Gossip算法中的节点没有Leader和Follower之分，它们都是对等的。
假定在一个有边界的网络中，每个状态机节点随机和其他状态机节点通讯，经过一段时间后在这个网络中的所有状态机节点的状态会达成一致。每个状态机节点都可能知道其他状态机节点或者只知道他邻近的节点，通信后他们的最终状态是一致的。从这个算法模型上看确实很想流言和流行病的传播模型。

我在举个实际例子来理解这个模型：

我们的每个状态机节点带有三元祖数据（key，value，version），我们能保证version是单向递增的。假设我们节点中有A，B，C，E，F五个节点构成环形结构，然后每个节点和相邻的节点每隔1秒通讯一次，把自己的数据（key，value，version）推个两个相邻节点，这两个相邻节点拿到数据后和自己的数据做比对，如果数据比自己新则更新自己的数据。这个模型最多2秒就能让整个集群数据收敛成一致。

上面的例子还有很多需要讨论的，第一每次同步的数据是全量还是增量，第二每次更新是节点主动推还是节点主动去拉，还会既有推也有拉。这两个问题会衍生初Gossip的类型和通讯方式。


## Gossip 类型

Gossip 有两种类型：

- Anti-Entropy（反熵）：以固定的概率传播所有的数据，可以理解为全量比对。
- Rumor-Mongering（谣言传播）：仅传播新到达的数据，可以理解为增量比对。

Anti-Entropy模式会让节点物理资源（网络和cpu）负担很重，Rumor-Mongering模式对节点资源负担相对较小，但是如何界定新数据变得比较困难，而且很难容错，无法保证一致性，所以反而Anti-Entropy有更大的价值，对于Anti-Entropy模式的优化研究更多。Anti-Entropy模式并不是真正的传送所有数据，而是变成如何追踪整个数据的变动，然后快速的找到数据的差异将差异数据传送。默克尔树（Merkle Tree）就是非常不错的一差异定位算法，有兴趣的可以去了解下。


## Gossip算法的通讯方式

gossip Node A 和 Node B有三种通信方式:

push: A节点将数据和版本号(key,value,version)推送给B节点，B节点更新A中比自己新的数据
pull：A仅将（key,version）推送给B，B将本地比A新的数据（Key,value,version）推送给A，A更新本地
push/pull：A仅将（key,version）推送给B，B将本地比A新的数据（Key,value,version）推送给A，A更新本地，A再将本地比B新的数据推送给B，B更新本地。

从上面的描述上看，push/pull的方式虽然通讯次数最多但是仅需要一个时间周期就能让A,B节点完全同步，Cassandra就是采用这个方式。


## Gossip算法的优点：

- 扩展性：网络可以允许节点的任意增加和减少，新增加的节点的状态最终会与其他节点一致。
- 容错性：网络中任何节点的宕机和重启都不会影响 Gossip 消息的传播，Gossip 算法具有天然的分布式系统容错特性。
- 去中心化：Gossip 算法不要求任何中心节点，所有节点都是对等的，任何一个节点无需知道整个网络状况，只要网络是连通的，任意一个节点就可以把消息散播到全网。


## Gossip算法的缺点：

- Gossip算法无法确定某个时刻所有状态机的状态是否一致。
- Gossip算法由于要经常和自己的相关节点通讯，因此可能早大量冗余的网络流量，甚至可能造成流量风暴。


## 总结

Gossip算法从他的特性来说应该是一种非常妙的算法，在非强一致性要求的领域非常实用，去中心话，同时又有天然的拓展性，顺带天然的故障检测属性。
在go生态在gossip的实现比较多，比较出门的有hashicorp实现的[memberlist](https://github.com/hashicorp/memberlist)。

## 参考资料
- [flowgossip](http://www.cs.cornell.edu/home/rvr/papers/flowgossip.pdf)