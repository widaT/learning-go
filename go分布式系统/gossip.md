# 分布式一致性算法——Gossip算法

前两节我们介绍了CA算法的典型——2PC，CP算法的典型——Raft，本小节要介绍AP算法的典型——Gossip算法。

Gossip算法有很多别名，流言算法，流行病算法等。

我们看下Gossip的算法模型：
假定在一个有边界的网络中，每个状态机节点随机和其他状态机节点通讯，经过一段时间后在这个网络中的所有状态机节点的状态会达成一致。每个状态机节点都可能知道其他状态机节点或者只知道他邻近的节点，通信后他们的最终状态是一致的。从这个算法模型上看确实很想流言和流行病的传播模型。






Gossip 类型

Gossip 有两种类型：

Anti-Entropy（反熵）：以固定的概率传播所有的数据
Rumor-Mongering（谣言传播）：仅传播新到达的数据
Anti-Entropy 是 SI model，节点只有两种状态，Suspective 和 Infective，叫做 simple epidemics。

Rumor-Mongering 是 SIR model，节点有三种状态，Suspective，Infective 和 Removed，叫做 complex epidemics。

其实，Anti-entropy 反熵是一个很奇怪的名词，之所以定义成这样，Jelasity 进行了解释，因为 entropy 是指混乱程度（disorder），而在这种模式下可以消除不同节点中数据的 disorder，因此 Anti-entropy 就是 anti-disorder。换句话说，它可以提高系统中节点之间的 similarity。

在 SI model 下，一个节点会把所有的数据都跟其他节点共享，以便消除节点之间数据的任何不一致，它可以保证最终、完全的一致。

由于在 SI model 下消息会不断反复的交换，因此消息数量是非常庞大的，无限制的（unbounded），这对一个系统来说是一个巨大的开销。



但是在 Rumor Mongering（SIR Model） 模型下，消息可以发送得更频繁，因为消息只包含最新 update，体积更小。而且，一个 Rumor 消息在某个时间点之后会被标记为 removed，并且不再被传播，因此，SIR model 下，系统有一定的概率会不一致。

而由于，SIR Model 下某个时间点之后消息不再传播，因此消息是有限的，系统开销小。









## Gossip算法的通讯方式

和Raft算法不同，Gossip算法中的节点没有Leader和Follower之分，它们都是对等的。

gossip Node A 和 Node B有三种通信方式:

push: A节点将数据(key,value,version)及对应的版本号推送给B节点，B节点更新A中比自己新的数据
pull：A仅将数据key,version推送给B，B将本地比A新的数据（Key,value,version）推送给A，A更新本地
push/pull：与pull类似，只是多了一步，A再将本地比B新的数据推送给B，B更新本地如果把两个节点数据同步一次定义为一个周期，则在一个周期内，push需通信1次，pull需2次，push/pull则需3次，从效果上来讲，push/pull最好，理论上一个周期内可以使两个节点完全一致。直观上也感觉，push/pull的收敛速度是最快的。

假设每个节点通信周期都能选择（感染）一个新节点，则Gossip算法退化为一个二分查找过程，每个周期构成一个平衡二叉树，收敛速度为O(n2 )，对应的时间开销则为O(logn )。这也是Gossip理论上最优的收敛速度。但在实际情况中最优收敛速度是很难达到的，假设某个节点在第i个周期被感染的概率为pi ,第i+1个周期被感染的概率为pi+1 ，则pull的方式:



而push为：



显然pull的收敛速度大于push，而每个节点在每个周期被感染的概率都是固定的p(0<p<1)，因此Gossip算法是基于p的平方收敛，也成为概率收敛，这在众多的一致性算法中是非常独特的。

个Gossip的节点的工作方式又分两种：

Anti-Entropy（反熵）：以固定的概率传播所有的数据
Rumor-Mongering（谣言传播）：仅传播新到达的数据
Anti-Entropy模式有完全的容错性，但有较大的网络、CPU负载；Rumor-Mongering模式有较小的网络、CPU负载，但必须为数据定义”最新“的边界，并且难以保证完全容错，对失败重启且超过”最新“期限的节点，无法保证最终一致性，或需要引入额外的机制处理不一致性。我们后续着重讨论Anti-Entropy模式的优化。



我们看下Gossip算法的优点：

- 扩展性：网络可以允许节点的任意增加和减少，新增加的节点的状态最终会与其他节点一致。
- 容错性：网络中任何节点的宕机和重启都不会影响 Gossip 消息的传播，Gossip 算法具有天然的分布式系统容错特性。
- 去中心化：Gossip 算法不要求任何中心节点，所有节点都可以是对等的，任何一个节点无需知道整个网络状况，只要网络是连通的，任意一个节点就可以把消息散播到全网。
- 最终一致性：Gossip 算法在网络中通过相互的通讯传播，可以达到指数级的传播速度，然后整个网络状态最终后收敛到一致。本质上将增加节点不会影响Gossip的传播性能，因此Gossip具有幂等性。

Gossip算法的缺点：

- Gossip算法无法确定某个时刻所有状态机的状态是否一致。
- Gossip算法由于要经常和自己的相关节点通讯，因此可能早大量冗余的网络流量，甚至可能造成流量风暴。


## Gossip golang实现

[memberlist](https://github.com/hashicorp/memberlist)



## 参考资料

- [flowgossip](http://www.cs.cornell.edu/home/rvr/papers/flowgossip.pdf)