# go WebSocket

## 什么是WebSocket

WebSocket是一种网络传输协议，可在单个TCP连接上进行全双工通信。其实看这定义会发现，跟其他Socket的协议没什么区别，都是全双工通信，叫WebSocket有何深意？
WebSocket使用了http和https的端口也就是80和443端口，同时第一次握手（handshaking）是基于http1.1。

## WebSocket使用场景

WebSocket主要还是使用来浏览器和后端交互上，特别是js和后端服务的交互。在WebSocket出现之前，浏览器js和后端难以建立长连接，后端主动通知前端没有有效途径，只能靠前端自己轮循后端接口实现，这样子的做法低效而且相当耗资源。WebSocket的出现正好可以解决这个痛点。


## WebSocket 握手过程
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

## WebSocket数据帧

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


###　实现“回音程序”
我们写一个“回音程序”（client给服务端发消息，服务端回复给你同样的消息）。

```go
package main
import (
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
	"log"
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
func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}
func main() {
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(addr, nil))
}
var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>点击 "打开" 创建websocket连接, 
点击 "发送" 将发送信息给服务端 ， 点击"关闭" 将关闭websocket连接.
<p>
<form>
<button id="open">打开</button>
<button id="close">关闭</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">发送</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
```
运行程序
```bash
$ go run main.go
```
然后用chrome浏览器器访问`http://localhost:8080/`你会看到页面，按照页面提示操作，体验WebSocket创建和发送消息。

### 给服务端发shell命令

WebSocket可以让浏览器后端见了一个全双工通道，本质上就是一个tcp连接，所以能做的事情应该很多，比如ssh over WebSocket。
这个demo程序展示客服端发shell 命令在服务端执行并返回执行结果。

```bash
$ tree
.
├── go.mod
├── go.sum
├── home.html
└── main.go
```


main.go内容：
```go
package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
	"github.com/gorilla/websocket"
)

var (
	addr    = flag.String("addr", "127.0.0.1:8080", "http service address")
	cmdPath string
)

const (
	writeWait = 10 * time.Second
	maxMessageSize = 8192
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	closeGracePeriod = 10 * time.Second
)

func pumpStdin(ws *websocket.Conn, w io.Writer) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = append(message, '\n')
		if _, err := w.Write(message); err != nil {
			break
		}
	}
}

func pumpStdout(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	defer func() {
	}()
	s := bufio.NewScanner(r)
	for s.Scan() {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}
	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}
func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}
var upgrader = websocket.Upgrader{}
func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()

	outr, outw, err := os.Pipe()
	if err != nil {
		internalError(ws, "stdout:", err)
		return
	}
	defer outr.Close()
	defer outw.Close()

	inr, inw, err := os.Pipe()
	if err != nil {
		internalError(ws, "stdin:", err)
		return
	}
	defer inr.Close()
	defer inw.Close()

	proc, err := os.StartProcess(cmdPath, flag.Args(), &os.ProcAttr{
		Files: []*os.File{inr, outw, outw},
	})
	if err != nil {
		internalError(ws, "start:", err)
		return
	}

	inr.Close()
	outw.Close()

	stdoutDone := make(chan struct{})
	go pumpStdout(ws, outr, stdoutDone)
	go ping(ws, stdoutDone)

	pumpStdin(ws, inw)
	inw.Close()
	if err := proc.Signal(os.Interrupt); err != nil {
		log.Println("inter:", err)
	}

	select {
	case <-stdoutDone:
	case <-time.After(time.Second):
		// A bigger bonk on the head.
		if err := proc.Signal(os.Kill); err != nil {
			log.Println("term:", err)
		}
		<-stdoutDone
	}
	if _, err := proc.Wait(); err != nil {
		log.Println("wait:", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	var err error
	cmdPath, err = exec.LookPath("sh")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
```
home.html内容

```html
<!DOCTYPE html>
<html lang="en">
<head>
	<title>Command Example</title>
	<script type="text/javascript">
        window.onload = function () {
            var conn;
            var msg = document.getElementById("msg");
            var log = document.getElementById("log");
            function appendLog(item) {
                var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
                log.appendChild(item);
                if (doScroll) {
                    log.scrollTop = log.scrollHeight - log.clientHeight;
                }
            }
            document.getElementById("form").onsubmit = function () {
                if (!conn) {
                    return false;
                }
                if (!msg.value) {
                    return false;
                }
                conn.send(msg.value);
                msg.value = "";
                return false;
            };
            if (window["WebSocket"]) {
                conn = new WebSocket("ws://" + document.location.host + "/ws");
                conn.onclose = function (evt) {
                    var item = document.createElement("div");
                    item.innerHTML = "<b>Connection closed.</b>";
                    appendLog(item);
                };
                conn.onmessage = function (evt) {
                    var messages = evt.data.split('\n');
                    for (var i = 0; i < messages.length; i++) {
                        var item = document.createElement("div");
                        item.innerText = messages[i];
                        appendLog(item);
                    }
                };
            } else {
                var item = document.createElement("div");
                item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
                appendLog(item);
            }
        };
	</script>
	<style type="text/css">
		html {
			overflow: hidden;
		}
		body {
			overflow: hidden;
			padding: 0;
			margin: 0;
			width: 100%;
			height: 100%;
			background: gray;
		}
		#log {
			background: white;
			margin: 0;
			padding: 0.5em 0.5em 0.5em 0.5em;
			position: absolute;
			top: 0.5em;
			left: 0.5em;
			right: 0.5em;
			bottom: 3em;
			overflow: auto;
		}
		#log pre {
			margin: 0;
		}
		#form {
			padding: 0 0.5em 0 0.5em;
			margin: 0;
			position: absolute;
			bottom: 1em;
			left: 0px;
			width: 100%;
			overflow: hidden;
		}
	</style>
</head>
<body>
<div id="log"></div>
<form id="form">
	<input type="submit" value="Send" />
	<input type="text" id="msg" size="64"/>
</form>
</body>
</html>
```

运行程序 
```bash
$ go run main.go
```

用chrome浏览器打开`http://localhost:8080/`，在最下发发送`ls`或者`netstat`等等命令，可以在显示区域看到执行结果。


## 总结

WebSocket一般不会用来做后端后后端的通信，经常来做web前端和后端的通信。本小节简要的介绍了WebSocket的工作原理和它的常见使用场景。以及用两个demo展示WebSocket的能力。其实还有很多场景可以去挖掘，感兴趣的同学可以去研究下。

## 参考资料

-  [ The WebSocket Protocol](https://tools.ietf.org/html/rfc6455)