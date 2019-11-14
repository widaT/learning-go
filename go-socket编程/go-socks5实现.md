# go socks5实现

## 什么是SOCKS

SOCKS是一种网络传输协议，主要用于Client端与外网Server端之间通讯的中间传递。SOCKS是"SOCKETS"的缩写，意外这一切socks都可以通过代理。
SOCKS是一种代理协议，相比于常见的HTTP代理，SOCKS的代理更加的底层，传统的HTTP代理会修改HTTP头，SOCKS则不会它只是做了数据的中转。

SOCKS的最新版本为SOCKS5。

## SOCKS5协议详解

SOCKS5的协议是基于二进制的，协议设计十分的简洁。

Client端和Server端第一次通讯的时候，Client端会想服务端发送版本认证信息，协议格式如下：
```
	+----+----------+----------+
	|VER | NMETHODS | METHODS  |
	+----+----------+----------+
	| 1  |    1     | 1 to 255 |
	+----+----------+----------+
```    
- VER是SOCKS版本，目前版本是0x05；
- NMETHODS是METHODS部分的长度；
- METHODS是Client端支持的认证方式列表，每个方法占1字节。当前的定义是：
  - 0x00 不需要认证
  - 0x01 GSSAPI
  - 0x02 用户名、密码认证
  - 0x03 - 0x7F由IANA分配（保留）
  - 0x80 - 0xFE为私人方法保留
  - 0xFF 无可接受的方法

Server端从Client端提供的方法中选择一个并通过以下消息通知Client端：

```
	+----+--------+
	|VER | METHOD |
	+----+--------+
	| 1  |   1    |
	+----+--------+
```    

- VER是SOCKS版本，目前是0x05；
- METHOD是服务端选中的方法。如果返回0xFF表示没有一个认证方法被选中，Client端需要关闭连接。

SOCKS5请求协议格式：

```
	+----+-----+-------+------+----------+----------+
	|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
```
- VER是SOCKS版本，目前是0x05；
- CMD是SOCK的命令码
    - 0x01表示CONNECT请求
    - 0x02表示BIND请求
    - 0x03表示UDP转发
- RSV 0x00，保留
- ATYP DST.ADDR类型
    - 0x01 IPv4地址，DST.ADDR部分4字节长度
    - 0x03 域名地址，DST.ADDR部分第一个字节为域名长度，剩余的内容为域名。
    - 0x04 IPv6地址，16个字节长度。
- DST.ADDR 目的地址
- DST.PORT 网络字节序表示的目的端口

Server端按以下格式返回给Client端：

```
	+----+-----+-------+------+----------+----------+
	|VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
```
- VER是SOCKS版本，目前是0x05；
- REP应答字段
    - 0x00 表示成功
    - 0x01 普通SOCKSServer端连接失败
    - 0x02 现有规则不允许连接
    - 0x03 网络不可达
    - 0x04 主机不可达
    - 0x05 连接被拒
    - 0x06 TTL超时域名地址
    - 0x07 不支持的命令
    - 0x08 不支持的地址类型
    - 0x09 - 0xFF未定义
- RSV 0x00，保留
- ATYP BND.ADDR类型
    - 0x01 IPv4地址，DST.ADDR部分4字节长度
    - 0x03 域名地址，DST.ADDR部分第一个字节为域名长度，剩余的内容为域名。
    - 0x04 IPv6地址，16个字节长度。
- BND.ADDR Server端绑定的地址
- BND.PORT 网络字节序表示的Server端绑定的端口

## 用go实现socks5

```golang
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)
const (
	IPV4ADDR = uint8(1) //ipv4地址
	DNADDR   = uint8(3) //域名地址
	IPV6ADDR = uint8(4) //地址
	CONNECTCMD = uint8(1)
	SUCCEEDED                     = uint8(0)
	NETWORKUNREACHABLE            = uint8(3)
	HOSTUNREACHABLE               = uint8(4)
	CONNECTIONREFUSED             = uint8(5)
	COMMANDNOTSUPPORTED           = uint8(7)
)

type Addr struct {
	Dn   string
	IP   net.IP
	Port int
}
func (a Addr) Addr() string {
	if 0 != len(a.IP) {
		return net.JoinHostPort(a.IP.String(), strconv.Itoa(a.Port))
	}
	return net.JoinHostPort(a.Dn, strconv.Itoa(a.Port))
}


type Command struct {
	Version      uint8
	Command      uint8
	RemoteAddr   *Addr
	DestAddr     *Addr
	RealDestAddr *Addr
	reader       io.Reader
}

func auth(conn io.ReadWriter) error {
	header := make([]byte,2)
	_,err:= io.ReadFull(conn,header)
	if err!=nil {
		return err
	}
	//igonre check version
	methods := make([]byte,int(header[1]))
	_,err = io.ReadFull(conn,methods)
	if err!=nil {
		return err
	}
	_,err = conn.Write([]byte{uint8(5), uint8(0)}) // 返回协议5，不需要认证
	return err
}


func readAddr(r io.Reader) (*Addr, error) {
	d := &Addr{}
	addrType := []byte{0}
	if _, err := r.Read(addrType); err != nil {
		return nil, err
	}
	switch addrType[0] {
	case IPV4ADDR:
		addr := make([]byte, 4)
		if _, err := io.ReadFull(r, addr); err != nil {
			return nil, err
		}
		d.IP = net.IP(addr)
	case IPV6ADDR:
		addr := make([]byte, 16)
		if _, err := io.ReadFull(r, addr); err != nil {
			return nil, err
		}
		d.IP = net.IP(addr)

	case DNADDR:
		if _, err := r.Read(addrType); err != nil {
			return nil, err
		}
		addrLen := int(addrType[0])
		DN := make([]byte, addrLen)
		if _, err := io.ReadFull(r, DN); err != nil {
			return nil, err
		}
		d.Dn = string(DN)

	default:
		return nil, errors.New("unkown addr type")
	}

	port := []byte{0, 0}
	if _, err := io.ReadFull(r, port); err != nil {
		return nil, err
	}
	d.Port = (int(port[0]) << 8) | int(port[1])
	return d, nil
}


func request(conn io.ReadWriter) (*Command, error)  {
	header := []byte{0, 0, 0}
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}
	//igonre check version
	dest, err := readAddr(conn)
	if err != nil {
		return nil, err
	}
	cmd := &Command{
		Version:  uint8(5),
		Command:  header[1],
		DestAddr: dest,
		reader:   conn,
	}

	return cmd, nil
}

func replyMsg(w io.Writer, resp uint8, addr *Addr) error {
	var addrType uint8
	var addrBody []byte
	var addrPort uint16
	switch {
	case addr == nil:
		addrType = IPV4ADDR
		addrBody = []byte{0, 0, 0, 0}
		addrPort = 0

	case addr.Dn != "":
		addrType = DNADDR
		addrBody = append([]byte{byte(len(addr.Dn))}, addr.Dn...)
		addrPort = uint16(addr.Port)

	case addr.IP.To4() != nil:
		addrType = IPV4ADDR
		addrBody = []byte(addr.IP.To4())
		addrPort = uint16(addr.Port)

	case addr.IP.To16() != nil:
		addrType = IPV6ADDR
		addrBody = []byte(addr.IP.To16())
		addrPort = uint16(addr.Port)

	default:
		return errors.New("format address error")
	}
	bodyLen := len(addrBody)
	msg := make([]byte, 6+bodyLen)
	msg[0] = uint8(5)
	msg[1] = resp
	msg[2] = 0 // RSV
	msg[3] = addrType
	copy(msg[4:], addrBody)
	msg[4+bodyLen] = byte(addrPort >> 8)
	msg[4+bodyLen+1] = byte(addrPort & 0xff)
	_, err := w.Write(msg)
	return err
}

func handleSocks5(conn io.ReadWriteCloser) error {
	if err := auth(conn);err !=nil {
		return err
	}
	cmd,err := request(conn)
	if err !=nil {
		return  err
	}
	fmt.Printf("%v",cmd.DestAddr)
	if err := handleCmd( cmd, conn); err != nil {
		return err
	}
	return nil
}

func handleCmd(cmd *Command, conn io.ReadWriteCloser) error {
	dest := cmd.DestAddr
	if dest.Dn != "" {
		addr, err := net.ResolveIPAddr("ip", dest.Dn)
		if err != nil {
			if err := replyMsg(conn, HOSTUNREACHABLE, nil); err != nil {
				return err
			}
			return err
		}
		dest.IP = addr.IP
	}

	cmd.RealDestAddr = cmd.DestAddr
	switch cmd.Command {
	case CONNECTCMD:
		return handleConn(conn, cmd)
	default:
		if err := replyMsg(conn, COMMANDNOTSUPPORTED, nil); err != nil {
			return err
		}
		return errors.New("Unsupported command")
	}
}

func handleConn(conn io.ReadWriteCloser, req *Command) error {
	target, err := net.Dial("tcp", req.RealDestAddr.Addr())
	if err != nil {
		msg := err.Error()
		resp := HOSTUNREACHABLE
		if strings.Contains(msg, "refused") {
			resp = CONNECTIONREFUSED
		} else if strings.Contains(msg, "network is unreachable") {
			resp = NETWORKUNREACHABLE
		}
		if err := replyMsg(conn, resp, nil); err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("Connect to %v failed: %v", req.DestAddr, err))
	}
	defer target.Close()

	local := target.LocalAddr().(*net.TCPAddr)
	bind := Addr{IP: local.IP, Port: local.Port}
	if err := replyMsg(conn, SUCCEEDED, &bind); err != nil {
		return err
	}
	go io.Copy(target, req.reader)
	io.Copy(conn, target)
	return nil
}
func main()  {
	l,err := net.Listen("tcp",":8090")
	if err !=nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		handle := func(conn net.Conn) {
			handleSocks5(conn)
		}
		go handle(conn)
	}
}
```

```bash
$ go run main.go
&{ 180.101.49.12 443} ## 使用curl出现
$ export ALL_PROXY=socks5://127.0.0.1:8090 # 设置终端代理
$ curl https://www.baidu.com
<!DOCTYPE html>
...
```
