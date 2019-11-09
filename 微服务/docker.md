# docker

## docker简介

docker是基于golang语言开发，基于linux kernel的cgroup和namespace，以及union Fs等技术对进程进行封装和隔离的轻量级容器。

<img src="../img/docker1.png" width = "40%" /><img src="../img/docker2.png" width = "40%" />

相比较传统的虚拟机（VM）,docker和其他docker容器贡献宿主机内核，所以运行一个docker容器中的程序相当快。而VM则占用很多程序逻辑以外的开销。

总的说来docker有以下的特点：

- 灵活：即使最复杂的应用程序也可以容器化。
- 轻量级：容器利用并共享主机内核，在系统资源方面比虚拟机更有效。
- 可移植：您可以在本地构建，部署到云并在任何地方运行。
- 松散耦合：容器是高度自给自足并封装的容器，使您可以在不破坏其他容器的情况下更换或升级它们。
- 可扩展：您可以在数据中心内增加并自动分发容器副本。
- 安全：容器将积极的约束和隔离应用于流程，而无需用户方面的任何配置。

## docker架构

docker采用经典的CS（客户端-服务器）架构。Docker客户端与Docker守护进程进行通讯，Docker守护进程完成了构建，运行和分发Docker容器的繁重工作。

<img src="../img/docker3.png" width = "90%" />


### Docker守护进程

docker守护进程（dockerd）侦听docker API请求并管理docker对象，例如镜像，容器，网络和数据卷。docker守护进程还可以与其他docker守护进程组成集群。

### Docker client

docker client（docker）是docker用户与docker交互的主要方式。使用诸如docker run等命令时，docker client会将这些命令发送到docker守护进程执行这些命令。Docker client还可以与多个docker守护进程通信。

### Docker仓库

docker仓库存储Docker镜像。Docker Hub是任何人都可以使用的Docker仓库，默认情况下docker在[Docker Hub](http://hub.docker.com)上查找镜像。用户也可以运行自己的私人镜像


### Docker对象

- docker镜像

    Docker 镜像是一个特殊的文件系统，除了提供容器运行时所需的程序、库、资源配置等文件外，还包含了一些为运行时准备的一些配置参数（ 如匿名卷、环境
变量、用户等） 。

- docker容器

    镜像（ Image） 和容器（ Container） 的关系，就像是面向对象程序设计中和实例一样，镜像是静态的定义，容器是镜像运行时的实体。容器可以创建、启动、停止、删除、暂停等。容器的实质是进程，但与直接在宿主执行的进程不同，容器进程运行于属于自己的独立的命名空间。因此容器可以拥有自己的 root 文件系统、自己的网络配置、自己的进程空间，甚至自己的用户 ID 空间。容器内的进程是运行在一个隔离的环境里，使用起来，就好像是在一个独立于宿主的系统下操作一样。这种特性使得容器封装的应用比直接在宿主运行更加安全。

## linux安装docker

 使用一键安装命令
 ```bash
 curl -sSL http://acs-public-mirror.oss-cn-hangzhou.aliyuncs.com/docker-engine/internet | sh -
```

## docker的使用

### 镜像管理

- 获取镜像 `docker pull [OPTIONS] NAME[:TAG|@DIGEST]`

  exp:
    - `docker pull ubuntu:17.10` 从docker hub 现在镜像
    - `docker pull http://abc.com:/ubuntu:17.10` 从其他仓库地址现在镜像 
- 镜像列表 `docker images [OPTIONS] [REPOSITORY[:TAG]]`

  exp:
    - `docker images ` 列出本地镜像，默认不显示中间镜像（intermediate image），列表包含了 仓库名 、 标签 、 镜像 ID 、 创建时间 以及 所占用的空间 。
        ```
        $ docker images
        micro-server          v1.0                bf580c3b368f        7 weeks ago         28.3MB
        wordpress             <none>              42a9bf5a6127        8 weeks ago         502MB #悬虚镜像
        ```
    - `docker image ls -f dangling=true`查看所以悬虚镜像
    - `docker image prune` 删除悬虚镜像，悬虚镜像已经没有作用可以随意删除
    - `docker images -a` 列出本地所有的镜像
     ```
     $ docker images -a
     REPOSITORY            TAG                 IMAGE ID            CREATED             SIZE
     grafana/grafana       6.0.1               ffd9c905f698        8 months ago        241MB
     <none>                <none>              f5690672aa36        12 months ago       133MB   #中间层级镜像
     ```
   - `docker image ls -f`过滤镜像
   - `docker image ls --format` 特定格式显示列表

- 删除镜像 `docker rmi [OPTIONS] IMAGE [IMAGE...]`
exp:
   - `docker rmi  ffd9c905f698` 删除镜像 
   - `docker image rm $(docker image ls -q redis)` 复合命令删除名字为redis的所有镜像

- 制作镜像
    - 修改后的容器保存成镜像`docker commit --author "xxx<xxx@xxx.com>" --message "nginx" webserver nginx:v2`这个方案很少用
    - 使用Dockerfile定制镜像
        exp:
        ```
        FROM alpine:latest
        WORKDIR /
        COPY  cmd/server/server /
        #EXPOSE 8000
        CMD ["./server"]
        ``
- 镜像仓库管理
    - `docker tag SOURCE_IMAGE[:TAG] TARGET_IMAGE[:TAG]` 设置镜像标签，这
       exp ：
       ```
        docker tag 860c279d2fec wida/nginx:v1
       ```
    - `docker push [OPTIONS] NAME[:TAG]`向远程镜像仓库推送标签的镜像
       exp：
       `docker push wida/nginx:v1` 这边是往docker hub推送

### 容器管理

- 查看容器列表 `docker ps [-a]` 不带-a的只列出正在运行的容器，带-a的列出所有容器
- 启动容器
    - 新建启动
    exp：
        - `docker run ubuntu:14.04 /bin/echo 'Hello world'` 运行容器
        - `docker run -t -i ubuntu:14.04 /bin/bash` 交互运行容器，-t让Docker分配一个伪终端并绑定到容器的标准输入上， -i则让容器的标准输入保持打开。
        - `docker run -id -p 8000:80 --name webserver nginx:v2` 使用-d后台运行
    - 启动停止的容器`docker start container`
- 停止容器 `docker stop container`
- 进入容器 
    - `docker attach container` 不建议使用
    - `docker exec container`  
        exp:docker exec  -i -t  nginx /bin/bash 让容器打开终端交互模式
- 删除容器
    - `docker rm container` 删除制定容器，首先得stop容器
    - `docker container prune` 删除所有停止容器

### Docker数据持久化

docker 容器内是不适合做数据存储的，通常我们会把数据存到容器外，主要有以下两种方式：
- 数据卷
    - `docker volume create myvol` 创建数据卷
    - `docker run -id  -v myvol:/data/app --name damov1 damo:v1` 挂载数据卷到容器/data/app
- 挂载主机目录
    - `docker run -id  -v /data/local/app:/data/app --name damov1 damo:v1` 挂载宿主机目录到容器目录/data/app

### Docker网络

- 网络端口映射 
    - `-P`是容器内部端口随机映射到主机的高端口。
        - `docker run -d -P --name damov1 damo:v1` 随机映射端口到容器的开放端口上
    - `-p`是容器内部端口绑定到指定的主机端口。
        - `docker run -id -p 8010:50051 --name damov1 damo:v1` 将容器的端口50051映射到宿主机8010上
        - `docker run -id -p 127.0.0.1::50051 --name damov1 damo:v1`将容器50051映射到localhost的随机端口上
        - `docker run -id -p 8010:50051/udp --name damov1 damo:v1`指定network类型，容器udp 50051端口映射到宿主机udp 8010端口。不写默认都是tcp。

- 容器互联
    两个容器要互相通讯需要如下步骤
    - 建立新的网络 `docker network create -d bridge mynet`
    - 运行容器1指定network`docker run -it --rm --name server1 --network mynet busybox sh` 
    - 运行容器2指定network和容器1相同`docker run -it --rm --name server2 --network mynet busybox sh`
      ```
      / # ping server1
        PING server1 (172.18.0.2): 56 data bytes
        64 bytes from 172.18.0.2: seq=0 ttl=64 time=0.167 ms
        64 bytes from 172.18.0.2: seq=1 ttl=64 time=0.151 ms
     ```

## 总结

本文只是介绍docker常用的一些命令，关于的docker的高级运用，建议大家看下docker官方的文档[Docker Documentation](https://docs.docker.com/)。另外docker本身到的容器编排工具docker-swarm将在容器编排小节一块介绍。

## 参考资料

- [Docker Documentation](https://docs.docker.com/)
- [《docker_practice》](https://github.com/yeasy/docker_practice)