# golang-gin-starter

⭐ 一款基于 Golang + Gin + JWT + Swagger + 自定义 Logger 的快速开发模版

## 项目启动

```bash
go run main.go --conf=config.toml
```

## 接口文档

swagger：
http://localhost:9180/swagger/index.html

## 文件目录说明

- constants: 配置等常量信息，配置信息初始化
- dao: 应用程序的数据访问层，负责与数据库或其他数据存储系统进行交互
- entity：接口输入输出结构体定义信息
- router: 接口层，web 请求入口，接收传入参数，执行相关逻辑，并返回对应信息
- service: 业务逻辑层，处理相关逻辑信息，并调用 dao 层，存取相关信息
- utils: 工具函数文件夹
- vendor: 保存项目所依赖的第三方包的源代码副本
- build.sh: 工程编译脚本，运行./build.sh 将对 golang 服务编译打包
- config.toml： 工程配置文件，配置端口、日志路径等相关信息
- go.mod: 定义和管理项目的模块
- go.sum: 存储项目的模块依赖项的校验和信息
- main.go：程序入口，Config 初始化、Logger 初始化和 Db Client 初始化, gin 服务在此处启动；
- start.sh: 服务启动脚本，运行 ./start.sh，web 服务将被启动, 运行前需要先编译工程
- stop.sh: 服务停止脚本，运行 ./stop.sh，web 服务将被停止

## 编译后启动

```bash
# 编译工程
./build.sh

# 启动服务
./start.sh

# 停止服务
./stop.sh
```

如有问题请联系
+ QQ：2657977449 
+ 微信：zhang2657977449
