# http2协议详解

http2并不是http1的替代产物，应该说http2是http1的拓展，http2同header头压缩，网络连接复用（net connection multiplexing）

http2没有改动 http 的应用语义。 http方法、状态代码、URI 和标头字段等核心概念一如往常。不过http2 修改了数据格式化（分帧）以及在客户端与服务器间传输的方式。


## 建立http2连接

客户端不知道服务端是否支持http2时候如何建立http2连接呢？
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
 经过这样一次完整的握手之后就知道运用层协议了，如果是http2则协议标志是长度为2的字节数组[0x68, 0x32]代表"h2"。
- 非tls下则是针对"http"启动HTTP2，客户端无法预知服务端是否支持http2的情况先先发起http1.1请求，请求header包含HTTP2的升级报头字段和h2c标识，例如
    ```
    GET /default.htm HTTP/1.1
    Host: server.example.com
    Connection: Upgrade, HTTP2-Settings
    Upgrade: h2c
    HTTP2-Settings: <base64url encoding of HTTP/2 SETTINGS payload>
    ```
    如果服务端不支持http2则返回
    ```
    HTTP/1.1 200 OK
    Content-Length: 243
    Content-Type: text/html
    ```
    如果服务端支持http2则返回一个http状态码为101(转换协议)响应http头来接受升级请求。在101空内容响应终止后，服务端可以开始发送HTTP/2帧。这些帧必须包含一个发起升级的请求的响应。
    ```
    HTTP/1.1 101 Switching Protocols
    Connection: Upgrade
    Upgrade: h2c
    ```


## HTTP/2 Connection Preface

在http2建立连接后，服务端和客服端会互相发送24个字节内容为`PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n` 的连接序言（Connection Preface），然后紧接着会是一个`http2 frame` 类型为`SETTINGS（后面为讲解frame）`，这个帧可能为空。

到了这边双方算是已经完成连接了，后续客户端和服务端之间可以马上交换数据帧。

## http2 Frame

http2 frame格式如下

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


## 参考资料
- [http2-spec](https://http2.github.io/http2-spec/#starting)
- [rfc7301](https://tools.ietf.org/html/rfc7301)