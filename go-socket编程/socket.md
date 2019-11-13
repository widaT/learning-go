# Socket（套接字）编程

## 什么是Socket

Socket 是对 TCP/IP 协议族的一种封装，是应用层与TCP/IP协议族通信的中间软件抽象层。
Socket 还可以认为是一种网络间不同计算机上的进程通信的一种方法，利用三元组（ip地址，协议，端口）就可以唯一标识网络中的进程，网络中的进程通信可以利用这个标志与其它进程进行交互。
Socket 起源于 Unix ，Unix/Linux 基本哲学之一就是“一切皆文件”，都可以用“打开(open) –> 读写(write/read) –> 关闭(close)”模式来进行操作。因此 Socket 也被处理为一种特殊的文件。

## Socket常见类型

### Datagram sockets

无连接Socket，使用用户数据报协议（UDP）。在Datagram sockets上发送或接收的每个数据包都经过单独寻址和路由。数据报Socket无法保证顺序和可靠性，因此从一台机器或进程发送到另一台机器或进程的多个数据包可能以任何顺序到达或根本不到达。

### Stream sockets

面向连接的Socket，使用传输控制协议（TCP），流控制传输协议（SCTP）或数据报拥塞控制协议（DCCP）。Stream sockets提供无记录边界的有序且唯一的数据流，并具有定义明确的机制来创建和销毁连接以及检测错误。Stream sockets可靠地按顺序传输数据。在Internet上，Stream sockets通常在TCP之上实现，以便应用程序可以使用TCP/IP协议在任何网络上运行。


## C/S模式中的Sockets

提供应用程序服务的计算机进程称为服务器，并在启动时创建处于侦听状态的Socket。这些Socket正在等待来自客户端程序的连接。

通过为每个客户端创建子进程并在子进程与客户端之间建立TCP连接，TCP服务器可以同时为多个客户端提供服务。为每个连接创建唯一的专用Socket。当与远程Socket建立Socket到Socket的虚拟连接或虚拟电路（也称为TCP 会话）时，它们处于建立状态，从而提供双工字节流。

服务器可以使用相同的本地端口号和本地IP地址创建多个同时建立的TCPSocket，每个Socket都映射到其自己的服务器子进程，为客户端进程服务。由于远程Socket地址（客户端IP地址和/或端口号）不同，因此操作系统将它们视为不同的Socket。即因为它们具有不同的Socket对元组。

UDP Socket无法处于已建立状态，因为UDP是无连接的。因此netstat不会显示UDP Socket的状态。UDP服务器不会为每个并发服务的客户端创建新的子进程，但是同一进程将通过同一Socket顺序处理来自所有远程客户端的传入数据包。这意味着UDP Socket不是由远程地址标识的，而是仅由本地地址标识的，尽管每个消息都具有关联的远程地址。


## 总结

本小节简要的介绍了Sockets的一些概念，介绍了Socket常见的两种类型。Socket编程技术涵盖的面和知识体现相当广泛，同样它的运用更是广泛，好不夸张的说互联网的技术都是基于Socket的。

## 参考资料

[wikipedia](https://en.wikipedia.org/wiki/Network_socket)