# FinalCheckmateServer
## 简介
此项目是我们一款名为 Checkmate 多人在线游戏的服务端。游戏 Checkmate 项目可以在我们 Github 上继续了解。此项目是基于另一个名为 ServerFramework 的项目开发的。 ServerFramework 是我们开发的服务端框架，然后此服务端项目就是基于 ServerFramework 这一个服务端框架进行的服务端开发。

## 项目结构介绍
### 此服务端项目中存在2个 Server
服务端上的 Server 是按照 ServerFramework 框架中提供的 Server 的接口实现的。本项目中实现了两个与游戏相关的 Server.
* ZoneServer - 游戏区域服务器
此服务器的主要功能是进行客户端连接，游戏房间的管理。代码位于`/ZoneServer/`文件夹中

* GameServer - 游戏内容服务端
此服务器的主要功能是提供与进行游戏相关的范围。比如，基于框架提供的帧同步功能，此服务器进行客户端之间的帧同步消息转发。  
同时，玩家在 ZoneServer 中创建了游戏房间后，可通过框架提供的进程间通讯的方法，以RPC的方式，调用 GameServer 的功能，开启一局游戏  
代码位于`/gameserver/`文件夹下

### 此服务端上保存游戏相关数据
比如，游戏中用户可选择的角色，这一类数据，保存在此服务端上，位于`/Roles/`文件夹下。客户端在开始游戏时，服务端会把玩家所选择的人物的 json 数据发送给对应客户端。客户端再根据此数据，在客户端构造角色。

##  服务端的构建方法
在 FinalCheckmateServer 文件夹目录中执行`go build`; 之后，会产生出一个名为 FinalCheckmateServer 的可执行文件，则 build 成功。

## 服务端运行方法
1. 可选参数: `-ipModel`, 默认值: 0, 可选值: 0, 1。0代表获取本机的公网IP，然后会将这个IP发给客户端，从而与客户端建立具体游戏的连接；1代表获取本机的局域网IP，同样也会将这个IP作为之后创建游戏时与客户端交流的远端IP地址.
2. 运行方法: `./FinalCheckmateServer -ipModel=0` (-ipModel可选值为0，1。见上一条解释)
3. 服务端启动成功的标志: 控制台输出: `Start MainLoop`, 服务端已启动.....

## 客户端和服务端框架，可参阅我们Github代码仓库
* 客户端代码仓库：https://github.com/qaqgame/Checkmate
* 服务端框架代码仓库：https://github.com/qaqgame/ServerFramework
