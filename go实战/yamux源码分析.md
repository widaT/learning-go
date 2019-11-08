# 深入理解yamux

yamux是golang连接多路复用（connection multiplexing）的一个库，想法来源于google的SPDY（也就是后来的http2）。yamux能用很小的代价在一个真实连接（net connection）上实现上千个Client-Server逻辑流。

## 基本概念

- session（会话）
  session用于包裹（wrap）可靠的有序连接（net connection）并将其多路复用为多个流（stream）。

- stream（流）
  在session中，stream代表一个client-server逻辑流。stream有唯一且自增（+2）的id，客服端向server端的stream id为奇数，服务端向客户端发送的stream为偶数，同事0值代表session。stream是逻辑概念，传输的数据是以帧的形态传输的。
  
- frame（帧）  
  帧是在session中真正传输的数据，帧有两部分header和body
  - header包含12个字节的数据，（也就算每次消息发送会参数额外12字节）。
    - Version (8 bits) 协议版本，目前总是为0
    - Type (8 bits)  帧消息类型，
        - 0x0（Data）数据传输
        - 0x1（Window Update） 用更新stream收消息recvWindow的大小。（注意这个时候length字段则为窗口的增量值）
        - 0x2（Ping）心跳，keep alives和RTT度量作用
        - 0x3（Go Away） 用于关闭会话
    - Flags (16 bits)  
        - 0x1 SYN : 新stream需要被创建
        - 0x2 ACK : 确认新stream开始
        - 0x4 FIN : 执行stream的半关闭
        - 0x8 RST : 立即重置stream
    - StreamID (32 bits) 流ID用于区分逻辑流
    - Length (32 bits)  body的长度或者type为Window Update时的delta值
   - body 是真实需要传输的数据，可能没有。

## 实现原理 

### yamux multiplexing 如何实现的
<img src="../img/yamux1.png" width = "70%" />

 从上图我们可以看出multiplexing的原理：传输过程中使用frame传输，每个frame都带有stream ID，在传输过程中stream相同stream的数据有先后顺序但可能不是连续的，接收端通过逻辑映射关系整合成有序的stream。

### stream的状态变迁

类似tcp连接，每一个stream都是有链接状态的。一个新的stream在创建的时候client会向server发送SYN信息(这边和tcp有个不一样的地方是，SYN发送后可以立即发送数据而不是等等对方ACK后再发)，server端接收到SYN信息后会回传ACK信息。close的时候会给对方发送FIN，对方接收到同样会回一个FIN。这个过程会伴随整个stream的状态变迁。
 <img src="../img/yamux2.png" width = "60%" />

上图还有一种 `streamReset` 状态没有呈现，server端Accept等待队列满的时候会发`flagRST`送给client的信息，client收到这个消息后会把流状态设置成`streamReset`这个时候流会停止。

### 流控制(Flow Control)

类似TCP的流控制，yamux也提供一种机制可以让发送端根据接收端接收能力控制发送的数据的大小。Tcp流控制的操作是接收端向发送端通知自己可以接收数据的大小，发送端会发送不超过这个限度的数据。这个大小限度就被称作窗口（window）大小。我们从stream状态迁移的图中看到了一个概念-window（窗口），就是和Tcp窗口类似的概念。

yamux的每个stream的初始窗口为256k，当然这个值是是可以配置修改的，在stream的SYN和ACK的消息交互中就带了窗口大小的协商。

窗口的大小由接收端决定的，接收端将自己可以接收的缓冲区大小通`typeWindowUpdate`类型的Header发送给发送端，发送端根据这个值调整自己发送数据的大小，如果发现是0就会阻塞发送。

## 源码分析

### 创建session
创建session只能通过 `Server(conn io.ReadWriteCloser, config *Config) `和 `func Client(conn io.ReadWriteCloser, config *Config) ` 这两个方法创建,本质上都调用了`newSession`的方法。我们具体看下`newSession`的方法。

```golang
func newSession(config *Config, conn io.ReadWriteCloser, client bool) *Session {
...
	s := &Session{
		config:     config,
		logger:     logger,
		conn:       conn,                  //真实连接（实际上io.ReadWriteCloser）
		bufRead:    bufio.NewReader(conn),           
		pings:      make(map[uint32]chan struct{}),
		streams:    make(map[uint32]*Stream),         //流映射
		inflight:   make(map[uint32]struct{}),
		synCh:      make(chan struct{}, config.AcceptBacklog), 
		acceptCh:   make(chan *Stream, config.AcceptBacklog), // 控制accept 队列长度
		sendCh:     make(chan sendReady, 64), //发送的队列长度
		recvDoneCh: make(chan struct{}),    //recvloop 终止信号
		shutdownCh: make(chan struct{}),     //session 关闭信号
	}
	if client {
		s.nextStreamID = 1    //client端的话初始stream Id是1
	} else {
		s.nextStreamID = 2    //client端的话初始stream Id是2
	}
	go s.recv()      //循环读取真实连接frame数据，然后分发到相应的message-type handler
	go s.send()      //循环通过真实连接发送frame数据
	if config.EnableKeepAlive {
		go s.keepalive() //通过`ping`心跳保持连接
	}
	return s
}
```

从上面的代码我们可以看出client和server唯一的区别是`nextStreamID`不一样，发送和接收数据的方式并没有区别。

`Session.recv() `实际调用的是`Session.recvLoop`方法
```golang
defer close(s.recvDoneCh)
	hdr := header(make([]byte, headerSize))
	for {
		// Read the header
		if _, err := io.ReadFull(s.bufRead, hdr); err != nil {  //读取 frame Header
			...省略错误处理代码
        }
        ...省略版本确认代码
		mt := hdr.MsgType()
		if mt < typeData || mt > typeGoAway {  //验证header type
			return ErrInvalidMsgType
		}

        if err := handlers[mt](s, hdr); err != nil {  //handle header type 
           ...省略代码
		}
	}
```

handlers 是一个全局变量，已经初始化的一个函数指针数组

```golang
handlers = []func(*Session, header) error{
		typeData:         (*Session).handleStreamMessage,
		typeWindowUpdate: (*Session).handleStreamMessage,
		typePing:         (*Session).handlePing,
		typeGoAway:       (*Session).handleGoAway,
	}
```

我们重点关注 `handleStreamMessage`方法，这个方法处理`typeWindowUpdate`和`typeData`这两个核心的消息类型。

```golang
func (s *Session) handleStreamMessage(hdr header) error {
	id := hdr.StreamID()
	flags := hdr.Flags()
	if flags&flagSYN == flagSYN {
		if err := s.incomingStream(id); err != nil { //如果是SYN信号，在接收端初始化stream，改变状态stream状态，通过 Session.acceptCh管道通知 接收端有新的stream
			return err
		}
	}
	// Get the stream
	s.streamLock.Lock()
	stream := s.streams[id]
	s.streamLock.Unlock()

     ...省略代码
	if hdr.MsgType() == typeWindowUpdate {
		if err := stream.incrSendWindow(hdr, flags); err != nil { //如果是typeWindowUpdate类型，这个调节本地发送窗口（sendWindow），触发`sendNotifyCh`可以发送更多数据
			...省略错误处理代码
			return err
		}
		return nil
	}
	if err := stream.readData(hdr, flags, s.bufRead); err != nil { //读取body
		...省略错误处理代码
		return err
	}
	return nil
}
```

总得来说 recv 的工作职责就是，读取每一个frame header，然后根据header type handle到不同的处理函数中。

`Session.send()`方法就相对简单些，职责是用底层真实链接发送frame给接收方。
```golang
func (s *Session) send() {
	for {
		select {
		case ready := <-s.sendCh:  //从缓冲管道（队列）获取frame
			if ready.Hdr != nil {
				sent := 0
				for sent < len(ready.Hdr) {
					n, err := s.conn.Write(ready.Hdr[sent:]) //先写header
					if err != nil {
					    ...省略错误处理代码
						return
					}
					sent += n
				}
			}
			if ready.Body != nil {  //如果有body写body
				_, err := io.Copy(s.conn, ready.Body)
				if err != nil {
					...省略错误处理代码
				}
			}
            ...省略版本确认代码
	}
}
```

# 参考文档

- [yamux](https://github.com/hashicorp/yamux/blob/master/spec.md)
- [go-libp2p 之 NewStream 深层阅读笔记](https://www.jianshu.com/p/14781d900501)
- [TCP流量控制与拥塞控制](https://www.jianshu.com/p/ad88e08e5dc8)