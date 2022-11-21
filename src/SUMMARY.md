# LEARNING-GO

# 前言

- [为什么选择go](./绪论.md)

# Go开发环境搭建
---
- [环境搭建](./go开发环境搭建/环境搭建.md)
    - [集成开发环境](./go开发环境搭建/集成开发工具.md)

# Go语言基础
---
- [第一个go程序](./go语言基础/第一个go程序.md)
    - [常量](./go语言基础/常量.md)
    - [变量](./go语言基础/变量.md)
    - [基本类型](./go语言基础/基本类型.md)
    - [字符串](./go语言基础/string.md)
    - [数组和切片](./go语言基础/数组和切片.md)
    - [map](./go语言基础/map.md)
    - [指针类型](./go语言基础/指针类型.md)
    - [运算](./go语言基础/运算.md)
    - [流程控制](./go语言基础/流程控制.md)
    - [函数](./go语言基础/函数.md)
    - [结构体和方法](./go语言基础/结构体和方法.md)
    - [package和可见性](./go语言基础/package和可见性.md)
    - [goroutine和channel](./go语言基础/goroutine和channel.md)
    - [interface](./go语言基础/interface.md)
    - [反射](./go语言基础/反射.md)
    - [错误处理](./go语言基础/错误处理.md)
    - [panic和recover](./go语言基础/panic和recover.md)
    
# Go标准库
---    
- [go标准库](./go标准库/go标准库概述.md)
    - [strings包](./go标准库/strings.md)
    - [bytes包](./go标准库/bytes.md)
    - [fmt格式化输入输出](./go标准库/fmt.md)
    - [文件读写](./go标准库/文件读写.md)
    - [time时间和日历](./go标准库/time.md)
    - [flag命令行参数解析](./go标准库/flag.md)
    - [json序列化](./go标准库/json.md)
    - [log程序日志](./go标准库/log.md)
    - [strconv字符串和其他基本类型转换](./go标准库/strconv.md)
    - [sort排序](./go标准库/sort.md)
    
# Go项目管理
---    
- [go项目管理](./go项目管理/go项目管理.md)
    - [go包依赖管理](./go项目管理/go-modules.md)
    - [go单元测试](./go项目管理/go-test.md)
    - [go cmd](./go项目管理/go命令.md)
    
# Go并发编程
---  
- [并发和并行](./go并发编程/并发和并行.md)
    - [原子操作](./go并发编程/原子操作.md)
    - [goroutine同步](./go并发编程/goroutine同步.md)
    - [条件变量](./go并发编程/条件变量.md)
    - [使用channel做goroutine同步](./go并发编程/使用channel做goroutine同步.md)
    - [sync once](./go并发编程/sync_once.md)
    - [并发安全map](./go并发编程/sync_map.md)
    - [sync pool](./go并发编程/sync_pool.md)
    - [同步拓展方案-semaphore](./go并发编程/semaphore.md)
    - [同步拓展方案-singleflight](./go并发编程/singleflight.md)
    - [同步拓展方案-errgroup](./go并发编程/errgroup.md)

# Go web服务
---  
- [go web服务](./web服务/go-web服务.md)
    - [http参数获取和http-client](./web服务/http参数获取和http-client.md)
    - [http请求参数校验](./web服务/http请求参数校验.md)
    - [web服务框架](./web服务/web服务框架.md)

# Go数据库编程
---  
- [go使用mysql](./go和数据库/go使用mysql.md)
    - [go使用mongodb](./go和数据库/go使用mongodb.md)
    - [go使用redis](./go和数据库/go使用redis.md)
    - [go使用es](./go和数据库/go使用es.md)

# Go内嵌kv数据库
---  
- [为什么需要内嵌数据库](./go内嵌kv数据库/README.md)
    - [boltdb](./go内嵌kv数据库/boltdb.md)
    - [badger](./go内嵌kv数据库/badger.md)
    - [goleveldb](./go内嵌kv数据库/goleveldb.md)

# Go socket编程
---  
- [socket编程](./go-socket编程/socket.md)
    - [go tcp udp服务](./go-socket编程/go-tcp-udp服务.md)
    - [go-socks5实现](./go-socket编程/go-socks5实现.md)
    - [go-websocket](./go-socket编程/go-websocket.md)
    - [深入理解connection-multiplexing](./go-socket编程/深入理解connection-multiplexing.md)

# Go微服务
---
- [微服务](./微服务/微服务.md)
    - [grpc和protobuf](./微服务/grpc和protobuf.md)
    - [docker](./微服务/docker.md)
    - [容器编排](./微服务/容器编排.md)
    - [go-micro微服务框架](./微服务/go-micro微服务框架.md)


# Go runtime
---
- [go runtime 简介](./go-runtime/runtime.md)
    - [go map](./go-runtime/go-map.md)
    - [go channel](./go-runtime/go-channel.md)
    - [go调度器](./go-runtime/scheduler.md)
    - [go内存分配器(未完)](./go-runtime/内存分配器.md)
    - [go垃圾回收器(未完)](./go-runtime/gc.md)

# Go和c语言
---
- [cgo简介](./cgo/cgo简介.md)
    - [cgo入门](./cgo/cgo入门.md)
    - [cgo使用场景](./cgo/cgo使用场景.md)
    - [指针和内存互访](./cgo/指针和内存互访.md)
    - [动态链接和静态链接](./cgo/动态链接和静态链接.md)

# Go汇编
---
- [go汇编简介(未完)](./go汇编/go汇编简介.md)
    - [go汇编基础(未完)](./go汇编/go汇编基础.md)
    - [函数(未完)](./go汇编/函数.md)
    - [go汇编使用场景(未完)](./go汇编/go汇编使用场景.md)
    - [go汇编笔记](./go汇编/go汇编笔记.md)

# Go分布式系统
---  
- [分布式系统简介](./go分布式系统/分布式系统.md)
    - [分布式一致性协议—2PC](./go分布式系统/2pc.md)
    - [分布式一致性协议—Raft](./go分布式系统/raft.md)
    - [分布式一致性协议—Gossip](./go分布式系统/gossip.md)
    - [分布式全局时间戳](./go分布式系统/全局时间戳.md)
    - [Etcd介绍](./go分布式系统/etcd.md)

# Go程序调试
---  
- [go调试器—delve](./go调试/go调试器.md)
    - [GODEBUG追踪调度器](./go调试/GODEBUG追踪调度器.md)
    - [GODEBUG追踪gc](./go调试/GODEBUG追踪gc.md)

# Go程序调优
---      
- [火焰图](./profiling/火焰图.md)

# Go实战
---  
- [go http2 开发](./go实战/go_http2开发.md)
    - [go语言结构体优雅初始化](./go实战/go语言结构体优雅初始化.md)
    - [go程序物理内存占用高的问题](./go实战/go程序物理内存占用高的问题.md)

# Demo代码
---
- [demo代码](./example_code.md)

# 关于我
---
- [关于我](./about.md)