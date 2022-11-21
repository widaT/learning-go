# go tcp、udp服务

golang标准库中的网络库非常之强大，那些在其他语言处理起来非常繁琐的socket代码，在golang中变得有点呆萌，而且却非常的高效。这是因为如此目前大行其道的云服务首先golang，很多区块链都是使用golang，分布式运用和容器编排软件通用使用golang。

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

### 解析文本协议解析

我们写一个解析redis 协议的demo
```golang
package main
import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)
var kvMap = make(map[string]string,10)
func parseCmd(buf []byte)([]string,  error){  //解析redis command 协议
	var cmd []string
	if buf[0] == '*' {
		for i:=1;i<len(buf);i++ {
			if buf[i] == '\n' && buf[i-1] =='\r'{
				count,_ := strconv.Atoi(string(buf[1:i-1]))
				for j:=0;j<count;j++ {
					i++
					if buf[i] != '$' {
						return nil,errors.New("error")
					}
					i++
					si:=i
					for ;i<len(buf);i++ {
						if buf[i] == '\n' && buf[i-1] =='\r' {
							size,_ := strconv.Atoi(string(buf[si:i-1]))
							cmd = append(cmd,string(buf[i+1:i+size+1]))
							i = i+size+2
							break
						}
					}
				}
			}
		}

	}
	return cmd,nil
}
func respString(msg string) []byte { //返回redis string
	b := bytes.Buffer{}
	b.Write( []byte{'+'})
	b.Write( []byte(msg))
	b.Write([]byte{'\r','\n'})
	return b.Bytes()
}
func respError(msg string) []byte{  //返回redis 错误信息
	b := bytes.Buffer{}
	b.Write( []byte{'-'})
	b.Write( []byte(msg))
	b.Write([]byte{'\r', '\n'})
	return b.Bytes()
}

func respNull()  []byte{      //返回redis Null
	return []byte{'$', '-', '1', '\r', '\n'}
}
func setKv(writer io.Writer,key,val string)  {
	kvMap[key] = val
	resp := respString("OK")
	writer.Write(resp)
}
func getV(writer io.Writer,key string)  {
	if v ,found := kvMap[key];found {
		resp := respString(v)
		writer.Write(resp)
	}else {
		writer.Write(respNull())
	}
}
func main() {
	addr := ":8088"
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
			defer conn.Close()
			buf := make([]byte,1024)
			n,_:=conn.Read(buf)
			cmd,err := parseCmd(buf[:n])
			if err !=nil {
				conn.Write(respError("COMMAND not supported"))
				return
			}
			switch strings.ToUpper(cmd[0]) {
			case "SET":
				setKv(conn,cmd[1],cmd[2])
			case "GET":
				getV(conn,cmd[1])
			default:
				conn.Write(respError("COMMAND not supported"))
			}
		}
		go handle(conn)
	}
}
```

```bash
$ go run main.go &
$ redis-cli -p 8088    ## 我们直接使用redis的cli
127.0.0.1:8088> set wida abx
OK
127.0.0.1:8088> get wida 
abx
127.0.0.1:8088> llen wida
(error) COMMAND not supported
127.0.0.1:8088> get amy
(nil)
127.0.0.1:8088>
```

这边只是简单实现了redis get和set的功能，有兴趣的同学可以拓展下实现redis的其他功能，这样子对你了解redis的底层原理有更深的认识。

### 解析二进制协议

我们先定义一下消息格式
```
/**
请求参数：
+-----+-------+------+----------+----------+----------+
| CMD | ARGS  |  L1  | STR1     |    Ln    | STRn     |
+-----+-------+------+----------+----------+----------+
|  1  |   1   |  1   | Variable |    2     | Variable |
+-----+-------+------+----------+----------+----------+
CMD 命令类型
ARG 参数个数
L1  参数一长度
STR1 参数与值
Ln  第n个参数长度
STRN 第n个参数值

返回格式
+-----+-------+----------+
| SUC |  LEN  |  BODY    |
+-----+-------+----------+
|  1  |   4   | Variable |
+-----+-------+----------+
SUC 是否成功
LEN BODY长度
BODY 消息体
*/
```

```golang
package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)
type Command struct {
	Cmd  uint8
	Args []string
}
func readCommand(r io.Reader) (*Command, error) {
	cmd := Command{}
	head := []byte{0, 0} //读取 CMD和ARGS
	_, err := io.ReadFull(r, head)
	if err != nil {
		return nil, err
	}
	cmd.Cmd = head[0]
	args := int(head[1])
	for i := 0; i < args; i++ { //循环读取 Ln STRn
		var length = []byte{0}
		_, err = r.Read(length)
		if err != nil {
			return nil, err
		}
		str := make([]byte, int(length[0]))
		_, err := io.ReadFull(r, str)
		if err != nil {
			return nil, err
		}
		cmd.Args = append(cmd.Args, string(str))
	}
	return &cmd, nil
}

func readResp(r io.Reader) (string, error) {
	suc := []byte{0}    //读取SUC 如果不成功就返回
	_, err := r.Read(suc)
	if err != nil {
		return "", err
	}
	if int(suc[0]) != 0 {
		return "", errors.New(fmt.Sprintf("resp errcode %v", suc[0]))
	}
	var bodyLen int32
	err = binary.Read(r, binary.BigEndian, &bodyLen) //大端读去长度 4字节
	if err != nil {
		return "", err
	}
	body := make([]byte, bodyLen)
	_, err = io.ReadFull(r, body) //读取body
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func main() {
	addr := ":8088"
	wg := sync.WaitGroup{}
	server := func() {
		wg.Add(1)
		defer wg.Done()
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
				defer conn.Close()
				for {
					cmd, err := readCommand(conn)
					if err != nil {
						break
					}
					fmt.Printf("server recv :%v\n",cmd)
					switch cmd.Cmd {
					case 1:
						//拼resp 字节
						conn.Write([]byte{uint8(0)})
						binary.Write(conn, binary.BigEndian, int32(10)) //大端写长长度 注意要和client约定好，当然可以用小端
						conn.Write([]byte("9876543210"))
					case 2:
						conn.Write([]byte{uint8(0)})
						binary.Write(conn, binary.BigEndian, int32(16))
						conn.Write([]byte("0000009876543210"))
					}
				}
			}
			go handle(conn)
		}
	}
	client := func() {
		wg.Add(1)
		defer wg.Done()
		conn, _ := net.Dial("tcp", addr)
		//拼CMD字节
		conn.Write([]byte{uint8(1), uint8(1), uint8(10)})
		conn.Write([]byte("0123456789"))
		ret, _ := readResp(conn)
		fmt.Println("client recv:",ret)
		//拼CMD字节
		conn.Write([]byte{uint8(2), uint8(1), uint8(16)})
		conn.Write([]byte("0123456789000000"))
		ret, _ = readResp(conn)
		fmt.Println("client recv:",ret)
	}
	go server()
	time.Sleep(1e9)
	go client()
	wg.Wait()
}
```

```bash
$ go run main.go
server recv :&{1 [0123456789]}
client recv: 9876543210
server recv :&{2 [0123456789000000]}
client recv: 0000009876543210
```

这边我们实现了二进制协议的解析，二进制协议的解析效率会比文本方式快，而且通常情况下会省带宽。


## 总结

本小节介绍了golang tcp和udp server端和client端的编写，我们还介绍了运用层协议文本协议和二进制协议的编写。本小节的demo会比较有启发性，大家可以发散思维做代码的拓展。