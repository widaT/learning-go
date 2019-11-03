# docker

## docker简介

docker是基于golang语言开发，基于linux kernel的cgroup和namespace，以及union Fs等技术对进程进行封装和隔离的轻量级容器。

<img src="../img/docker1.png" width = "40%" />
<img src="../img/docker2.png" width = "40%" />

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


### Docker 对象

- docker镜像

    Docker 镜像是一个特殊的文件系统，除了提供容器运行时所需的程序、库、资源配置等文件外，还包含了一些为运行时准备的一些配置参数（ 如匿名卷、环境
变量、用户等） 。

- docker容器

    镜像（ Image） 和容器（ Container） 的关系，就像是面向对象程序设计中和实例一样，镜像是静态的定义，容器是镜像运行时的实体。容器可以创建、启动、停止、删除、暂停等。容器的实质是进程，但与直接在宿主执行的进程不同，容器进程运行于属于自己的独立的命名空间。因此容器可以拥有自己的 root 文件系统、自己的网络配置、自己的进程空间，甚至自己的用户 ID 空间。容器内的进程是运行在一个隔离的环境里，使用起来，就好像是在一个独立于宿主的系统下操作一样。这种特性使得容器封装的应用比直接在宿主运行更加安全。

### linux安装docker

 使用一键安装命令
 ```bash
 curl -sSL http://acs-public-mirror.oss-cn-hangzhou.aliyuncs.com/docker-engine/internet | sh -
```

### docker的使用

```bash
启动docker
sudo service docker start

//守护进程的
docker run -id -p 8000:80 --name webserver nginx:v2 

//进入已经运行的容器内部
docker exec -it webserver bash

//commit 制作镜像
docker commit --author "xxx<xxx@xxx.com>" --message "nginx" webserver nginx:v2

//删除一个已经停止的容器
docker rm webserver

//删除镜像
docker rmi imageid

//删除所以已经停止运行的容器
docker rm $(docker ps -a -q)

//用dockerfile 构建镜像
docker build -t nginx:v3 .

//删除悬虚镜像
docker image prune

//数据卷
docker run -id -p 8010:50051 -v /path:/container_path/:rw --name demo demo:v1

//镜像过滤
docker image ls -f since=mongo:3.2

//docker push 到docker hub
//打tag
docker tag ubuntu:17.10 username/ubuntu:17.10

//push 
docker push username/ubuntu:17.10

//search
docker search username
```