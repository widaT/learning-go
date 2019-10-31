# docker

## 什么是docker

docker是基于golang语言开发，基于linux kernel的cgroup和namespace，以及union Fs等技术对进程进行封装和隔离的轻量级容器。

<img src="../img/docker1.png" width = "40%" />
<img src="../img/docker2.png" width = "40%" />

相比较传统的虚拟机（VM）,docker和其他docker容器贡献宿主机内核，所以运行一个docker容器中的程序相当快。而VM则占用很多程序逻辑以外的开销。

总的说来docker有如下的特点：

- 灵活：即使最复杂的应用程序也可以容器化。
- 轻量级：容器利用并共享主机内核，在系统资源方面比虚拟机更有效。
- 可移植：您可以在本地构建，部署到云并在任何地方运行。
- 松散耦合：容器是高度自给自足并封装的容器，使您可以在不破坏其他容器的情况下更换或升级它们。
- 可扩展：您可以在数据中心内增加并自动分发容器副本。
- 安全：容器将积极的约束和隔离应用于流程，而无需用户方面的任何配置。

## docker 架构

docker采用经典的CS（客户端-服务器）架构。Docker 客户端与Docker 守护进程进行通讯，Docker守护进程完成了构建，运行和分发Docker容器的繁重工作。

<img src="../img/docker3.png" width = "90%" />


### Docker守护进程

docker守护进程（dockerd）侦听docker API请求并管理docker对象，例如镜像，容器，网络和数据卷。docker守护进程还可以与其他docker守护进程组成集群。

### Docker client

docker client（docker）是docker用户与docker交互的主要方式。使用诸如docker run等命令时，docker client会将这些命令发送到docker守护进程执行这些命令。Docker client还可以与多个docker守护进程通信。

### Docker仓库

docker仓库存储Docker镜像。Docker Hub是任何人都可以使用的Docker仓库，默认情况下docker在[Docker Hub](http://hub.docker.com)上查找镜像。用户也可以运行自己的私人镜像。使用docker pull或docker run命令时，所需的镜像将从配置的Docker仓库中提取。使用该docker push命令时，会将镜像推送到配置的Docker仓库。


### Docker 对象

- docker镜像
    docker镜像是一种特殊的文件系统，
- docker容器