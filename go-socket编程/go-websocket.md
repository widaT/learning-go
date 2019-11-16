# go WebSocket

## 什么是WebSocket

WebSocket是一种网络传输协议，可在单个TCP连接上进行全双工通信。其实看这定义会发现，跟其他Socket的协议没什么区别，都是全双工通信，叫WebSocket有何深意？
WebSocket使用了http和https的端口也就是80和443端口，同时第一次握手（handshaking）是基于http1.1。

## WebSocket使用场景

WebSocket主要还是使用来浏览器和后端交互上，特别是js和后端服务的交互。在WebSocket出现之前，浏览器js和后端难以建立长连接，后端主动通知前端没有有效途径，只能靠前端自己轮循后端接口实现，低效而且耗资源。WebSocket的出现正好可以解决这个痛点。


### WebSocket 握手过程
一个典型的Websocket握手请求如下：

客户端请求
```
GET / HTTP/1.1
Upgrade: websocket
Connection: Upgrade
Host: example.com
Origin: http://example.com
Sec-WebSocket-Key: sN9cRrP/n9NdMgdcy2VJFQ==
Sec-WebSocket-Version: 13
```

服务器回应
```
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: fFBooB7FAkLlXgRSz0BT3v4hq5s=
Sec-WebSocket-Location: ws://example.com/
```

字段说明

- Connection必须设置Upgrade，表示客户端希望连接升级。
- Upgrade字段必须设置Websocket，表示希望升级到Websocket协议。
- Sec-WebSocket-Key是随机的字符串，服务器端会用这些数据来构造出一个SHA-1的信息摘要。把“Sec-WebSocket-Key”加上一个特殊字符串“258EAFA5-E914-47DA-95CA-C5AB0DC85B11”，然后计算SHA-1摘要，之后进行BASE-64编码，将结果做为“Sec-WebSocket-Accept”头的值，返回给客户端。如此操作，可以尽量避免普通HTTP请求被误认为Websocket协议。
- Sec-WebSocket-Version 表示支持的Websocket版本。RFC6455要求使用的版本是13，之前草案的版本均应当弃用。
- Origin字段是可选的，通常用来表示在浏览器中发起此Websocket连接所在的页面，类似于Referer。但是，与Referer不同的是，Origin只包含了协议和主机名称。

### WebSocket数据帧

WebSocket在握手之后不再使用文本协议，而是采用二进制协议。它的数据帧格式如下：
```
      0                   1                   2                   3
      0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
     +-+-+-+-+-------+-+-------------+-------------------------------+
     |F|R|R|R| opcode|M| Payload len |    Extended payload length    |
     |I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
     |N|V|V|V|       |S|             |   (if payload len==126/127)   |
     | |1|2|3|       |K|             |                               |
     +-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
     |     Extended payload length continued, if payload len == 127  |
     + - - - - - - - - - - - - - - - +-------------------------------+
     |                               |Masking-key, if MASK set to 1  |
     +-------------------------------+-------------------------------+
     | Masking-key (continued)       |          Payload Data         |
     +-------------------------------- - - - - - - - - - - - - - - - +
     :                     Payload Data continued ...                :
     + - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
     |                     Payload Data continued ...                |
     +---------------------------------------------------------------+
```

稍微有点复杂，各个字段的定义可以参考[RFC6455 5.2](https://tools.ietf.org/html/rfc6455#section-5.2)。这边不再赘述。

## go实现Websocket

golang的标准库就支持Websocket，但是我们推荐特性支持更多，而且性能更强的[gorilla](https://github.com/gorilla/websocket)。

我们写一个“回音程序”（client给服务端发消息，服务端回复给你同样的消息）。

```go
package main
import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var addr = "localhost:8080"
var upgrader = websocket.Upgrader{}
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for { //这边跟其他 和其他socket编程没什么区别，循环读取数消息
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("server recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func server() {
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func client() {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			log.Printf("client recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.Format("2006-01-02 15:04:05")))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func main()  {
	wg:= sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		server()
	}()

	time.Sleep(3e9)
	go func() {
		defer wg.Done()
		client()
	}()
	wg.Wait()
}
```

```bash
$ go run main.go
2019/11/16 17:03:58 connecting to ws://localhost:8080/echo
2019/11/16 17:03:59 server recv: 2019-11-16 17:03:59
2019/11/16 17:03:59 client recv: 2019-11-16 17:03:59
2019/11/16 17:04:00 server recv: 2019-11-16 17:04:00
```


## 总结

WebSocket一般不会用来做后端后后端的通信，经常来做web前端和后端的通信。本小节简要的介绍了WebSocket的工作原理和它的常见使用场景。WebSocket可以让浏览器后端见了一个全双工通道，本质上就是一个tcp连接，所以能做的事情应该很多，比如ssh over WebSocket等等。对前端感兴趣的可以深入研究下。

## 参考资料

-  [ The WebSocket Protocol](https://tools.ietf.org/html/rfc6455)