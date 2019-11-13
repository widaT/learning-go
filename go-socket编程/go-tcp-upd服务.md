# go tcp udp服务

golang标准库中的网络库非常之强大，那些在其他语言处理起来非常繁琐的socket代码，在golang中变得有点呆萌，而且却非常的高效。这是因为如此目前大行其道的云服务首先golang，很多区块链都是使用golang，分布式运用和容器编排软件通用使用golang。是不是有点迫不及待了，我们的后面的demo代码会尽量写得狂野点，enjoy！！！。


## Tcp服务

我们先写一个tcp C/S通讯模型的demo：
```golang
package main
import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)
func main() {
	addr :=":8088"
	wg := sync.WaitGroup{}
	server := func (){
		wg.Add(1)
		defer wg.Done()
		l,err := net.Listen("tcp",addr)
		if err !=nil {
			log.Fatal(err)
		}
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}
			handle := func(conn net.Conn) {
				buff := make([]byte,1024)
				n,_:=conn.Read(buff)
				fmt.Println(string(buff[:n]))
				conn.Write([]byte("nice to meet you too"))
			}
			go handle(conn)
		}
	}
	client := func(id string) {
		wg.Add(1)
		defer wg.Done()
		conn ,err:=	net.Dial("tcp",addr)
		if err !=nil {
			log.Fatal(err)
		}
		conn.Write([]byte("nice to meet you"))
		buff := make([]byte,1024)
		n,_:=conn.Read(buff)
		fmt.Println(id + " recv：" + string(buff[:n]))
	}

	go server()         //启动服务端
	time.Sleep(1e9)  //这边停1s等服务端启动
	go client("client1") //启动客服端1
	go client("client2") //启动客服端2
	wg.Wait()
}
```

```bash
$ go run main.go
nice to meet you
nice to meet you
client2 recv：nice to meet you too
client1 recv：nice to meet you to
```

我们用了51行代码，没有使用第三方库实现了tcp server和2个client的通讯。


## Udp服务
```golang
package main
import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)
func main() {
	addr :=":8088"
	wg := sync.WaitGroup{}
	server := func (){
		wg.Add(1)
		defer wg.Done()
		uAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			log.Fatal(err)
		}
		l,err := net.ListenUDP("udp",uAddr)
		if err !=nil {
			log.Fatal(err)
		}
		for {
			data := make([]byte,1024)
			n,rAddr,err :=l.ReadFrom(data)
			if err !=nil {
				log.Fatal(err)
			}
			fmt.Println(string(data[:n]))
			l.WriteTo([]byte("nice to meet you too"),rAddr)
		}
	}
	client := func(id string) {
		wg.Add(1)
		defer wg.Done()
		conn ,err:=	net.Dial("udp",addr)   //和tcp client代码相比，这边仅仅改了network类型
		if err !=nil {
			log.Fatal(err)
		}
		conn.Write([]byte("nice to meet you"))
		buff := make([]byte,1024)
		n,_:=conn.Read(buff)
		fmt.Println(id + " recv：" + string(buff[:n]))
	}

	go server()         //启动服务端
	time.Sleep(1e9)  //这边停1s等服务端启动
	go client("client1") //启动客服端1
	go client("client2") //启动客服端2
	wg.Wait()
}
```

```bash
$ go run main.go
nice to meet you
nice to meet you
client2 recv：nice to meet you too
client1 recv：nice to meet you to
```
在Server端代码稍微和tcp的方式有点不一样，client端仅仅把`net.Dial("tcp",addr) `改成`net.Dial("udp",addr)`同样非常简单。

## 运用层协议

网络协议是为计算机网络中进行数据交换而建立的规则、标准或约定的集合。通俗的说就是，双方约定的一种都能懂的数据格式。

应用层协议(application layer protocol)定义了运行在不同端系统上的应用程序进程如何相互传递报文。很多运用层协议是基于tcp的，因为可靠性对运用来说非常重要。我们常见的http，smtp（简单邮件传送协议），ftp都是基于tcp。比较少数的像dns（Domain Name System）是基于udp。

根据协议的序列化方式还可以分成 二进制协议和文本协议。

文本协议：一般是由一串ACSII字符组成，常见的文本协议有http1.0，[redis通讯协议（RESP）](https://redis.io/topics/protocol)。

二进制协议：有字节流数据组成，通常包括消息头(header)和消息体(body)，消息头的长度固定，并且定义了消息体的长度。

## 自定义协议


### 自定义文本协议解析
```golang

package main

import (
	"fmt"
	"strconv"
)

// *2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n

func parse(buf []byte)  {
	i:=0
	start:=0
	isNum := false
	for i < len(buf){
		switch  buf[i] {
			case '*':
				start = i+1
				isNum = true
		    case '\n':
		    	part := buf[start:i-1]
		    	if isNum {
		    		strconv.Atoi(string(part))
				}
		    	isNum = false
		    	fmt.Println(string(part))
				start = i+1
			case '$':
				start = i+1
				isNum = true
			default:
		}
		i++

	}
}


func main() {

	parse([]byte("*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n"))

	/*addr := ":8088"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		handle := func(conn net.Conn) {
			buf := make([]byte,1024)

			n,_:=conn.Read(buf)

			buf =buf[:n]





		}
		go handle(conn)
	}*/

}

```

### 自定义二进制协议
