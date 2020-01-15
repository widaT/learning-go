# Http2协议详解

Http2并不是http1的替代产物，应该说Http2是http1的拓展。Http2通过header头压缩，网络连接复用（network connection multiplexing）,server端推送，流优先级等手段来减少交互延时。

Http2没有改动 http 的应用语义。 http方法、状态代码、URI 和标头字段等核心概念一如往常。不过Http2 修改了数据格式化（分帧）以及在客户端与服务器间传输的方式。


## 建立Http2连接

客户端不知道服务端是否支持Http2时候如何建立Http2连接呢？
- 在有tls的情况下，会使用TLS的应用层协议协商扩展 [TLSALPN（Application-Layer Protocol Negotiation）]来协商tls后的通讯协议。下面展示一个完整的tls ALPN的握手过程。
    ```
     Client                                              Server

   ClientHello                     -------->       ServerHello
     (ALPN extension &                               (ALPN extension &
      list of protocols)                              selected protocol)
                                                   Certificate*
                                                   ServerKeyExchange*
                                                   CertificateRequest*
                                   <--------       ServerHelloDone
   Certificate*
   ClientKeyExchange
   CertificateVerify*
   [ChangeCipherSpec]
   Finished                        -------->
                                                   [ChangeCipherSpec]
                                   <--------       Finished
   Application Data                <------->       Application Data
    ```
 经过这样一次完整的握手之后就知道运用层协议了，如果是Http2则协议标志是长度为2的字节数组[0x68, 0x32]代表"h2"。
- 非tls下则是针对"http"启动Http2，客户端无法预知服务端是否支持Http2的情况先先发起http1.1请求，请求header包含Http2的升级报头字段和h2c标识，例如
    ```
    GET /default.htm HTTP/1.1
    Host: server.example.com
    Connection: Upgrade, Http2-Settings
    Upgrade: h2c
    Http2-Settings: <base64url encoding of HTTP/2 SETTINGS payload>
    ```
    如果服务端不支持Http2则返回
    ```
    HTTP/1.1 200 OK
    Content-Length: 243
    Content-Type: text/html
    ```
    如果服务端支持Http2则返回一个http状态码为101(转换协议)响应http头来接受升级请求。在101空内容响应终止后，服务端可以开始发送HTTP/2帧。这些帧必须包含一个发起升级的请求的响应。
    ```
    HTTP/1.1 101 Switching Protocols
    Connection: Upgrade
    Upgrade: h2c
    ```


## HTTP/2 Connection Preface

在Http2建立连接后，服务端和客服端会互相发送24个字节内容为`PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n` 的连接序言（Connection Preface），然后紧接着会是一个`Http2 frame` 类型为`SETTINGS（后面为讲解frame）`，这个帧可能为空。

到了这边双方算是已经完成连接了，后续客户端和服务端之间可以马上交换数据帧。

## Http2 Frame

Http2的所有数据交互是通过Http2 Frame实现的，Http2 frame是二级制实现的，相比于http1文本的实现，二进制在读取效率上有明显优势。

![](../img/Http2-1.svg)

新的二进制分帧机制改变了客户端与服务器之间交换数据的方式。 为了说明这个过程，我们需要了解 HTTP/2 的三个概念：

- Stream(数据流)：已建立的连接内的双向字节流，可以承载一条或多条消息。
- Message(消息)：与逻辑请求或响应消息对应的完整的一系列Frame。
- Frame（帧）：Http2通信的最小单位，每个Frame都包含Frame头，至少也会标识出当前Frame所属的Stream。

这些概念的关系总结如下：

- 所有通信都在一个TCP连接上完成，此连接可以承载任意数量的双向数据流。
- 每个数据流都有一个唯一的标识符和可选的优先级信息，用于承载双向消息。
- 每条消息都是一条逻辑HTTP消息（例如请求或响应），包含一个或多个Frame。
- Frame是最小的通信单位，承载着特定类型的数据，例如HTTP标头、消息负载等等。 来自不同数据流的Frame可以交错发送，然后再根据每个帧头的数据流标识符重新组装。

![](../img/http2-2.svg)

### 网络连接复用（Network Connection Multiplexing） 
在HTTP1中，如果客户端要想发起多个并行请求以提升性能，则必须使用多个TCP连接。 这是HTTP交付模型的直接结果，该模型可以保证每个连接每次只交付一个响应。更糟糕的是，这种模型也会导致队首阻塞，从而造成底层 TCP 连接的效率低下。

HTTP/2 中新的二进制分帧层突破了这些限制，实现了完整的网络连接复用：客户端和服务器可以将HTTP消息分解为互不依赖的帧，然后交错发送，最后再在另一端把它们重新组装起来。

![](../img/http2-3.svg)

关于Network Connection Multiplexing有专门的小节介绍[深入理解connection multiplexin](./深入理解connection-multiplexin.md)这边不再赘述。

### Http2 Frame格式

```
 +-----------------------------------------------+
 |                 Length (24)                   |
 +---------------+---------------+---------------+
 |   Type (8)    |   Flags (8)   |
 +-+-------------+---------------+-------------------------------+
 |R|                 Stream Identifier (31)                      |
 +=+=============================================================+
 |                   Frame Payload (0...)                      ...
 +---------------------------------------------------------------+
 ```
 总共9个字节的frame header头和长度为header中Length字段值的payload组成。

Header字段说明：

|字段 | 长度(bit) |  说明|  
|:-|-:|-:|
|Length | 24 |  payload的长度|
|Type | 8 | frame type（帧类型） |
|Flags | 8 | 对确定的帧类型赋予特定的语义, 否则发送时必须忽略(设置为0x0). |
|R | 1 | 预留字段,尚未定义语义. 发送和接收必须忽略(0x0).|
|Stream Identifier | 31  | 流标识 |


## 参考资料

- [http2-spec](https://http2.github.io/http2-spec/#starting)
- [rfc7301](https://tools.ietf.org/html/rfc7301)
- [http2](https://developers.google.com/web/fundamentals/performance/http2)