
<p align="center" style="font-size: 40px;font-weight: bold">
	Devman
</p>

<p align="center">
    <a target="_blank" href="">
        <img src='https://img.shields.io/badge/Golang-1.21.1-green.svg' alt='golang'/>
    </a>
    <a target="_blank" href="">
        <img src='https://img.shields.io/badge/Sqlite-3.0-orange.svg' alt='sqlite'/>
    </a>
    <a target="_blank" href="">
        <img src='https://img.shields.io/badge/Layui-2.8-yellow.svg' alt='sqlite'/>
    </a>
    <a target="_blank" href="">
        <img src='https://img.shields.io/badge/License-apache2.0-blue.svg' alt='License'/>
    </a>
</p>

<p align="center">
	<strong>🚀简单易用的库表结构全局展示软件</strong>
</p>

## 开篇
本来想用一些花里胡哨的名字, 如下
    比如mytool, 参考mycat, mysql之类, 还有本人非常喜欢的utool, hutool之类.
    又比如fatcat, 本人有一个胖猫, 通体白色, 爱称白胖, 从不粘人, 喜静不喜动, 也是它肉嘟嘟的原因. 就像本项目的初衷一样, 
它仅仅是一个工具而已, 轻巧, 不会自己加戏来吸引你的注意力. 你爱用就用, 不爱用就丢在一旁. 然而fatcat有大亨的意思, 很显然违背了初衷
因为本项目是以简洁轻快为宗旨开发的, 也是选用golang开发的原因之一.
最终都放弃了, 就像它提供的页面和功能都是如此直接了当, 没有任何套路. 也将简洁之风一以贯之.
所以就叫 `devman` 吧.


## 背景

### 元数据全局展示

场景一: 在做技术设计的时候, 经常需要查看表结构, 每次都要打开Navicat图形工具查看, 操作相对繁琐. 

场景二: 非开发人员需要查看表结构信息, 比如数据分析人员需要了解业务部门的底层数据关系, 在不需要授权的情况下, 可以直接查看库表结构.  



## 项目介绍
### 功能
为了解决上面出现的问题, 所以本项目开发了该功能
#### 1 元数据

![](www/static/img/devman01.png)

- 数据库表结构全局展示
- 支持单表分享
- 支持查看DDL建库脚本


## 开始

### 代码


``` shell
├── config  // 多环境配置文件
│   ├── boot.yml
│   ├── config-dev.yml
│   ├── config-prod.yml
│   └── config-test.yml
├── main.go // 程序入口
├── script // 脚本文件
├── src  // 后端代码
│   ├── common
│   ├── controllers // 控制层(mvc)结构
│   ├── middlewares // 中间件
│   ├── persistence  // 持久化层
│   │   ├── dbconn.go
│   │   └── model
│   ├── routers // 路由
└── www // 前端代码
    ├── html  // 模板文件
    │   ├── common
    │   │   └── page_footer.html
    │   └── index
    │       └── index.html
    └── static  // 静态资源
        ├── css
        ├── img
        └── js

```

### 部署


#### 方式一、本地运行
- 修改配置
config-dev.yml
```yaml
mysqls:  # db类型，后续支持pg、oracle等主流db
  - env: 生产环境   # 一级菜单 
    db: my_erpdb  # db名称  及二级菜单
    enable: true  # 是否开启
    host: localhost:3306  # db ip:端口
    user: root # db用户名
    password: 123456  # db密码
```

- 执行命令
```shell
# clone代码
go mod init devman
go mod tidy
go build
./devman
```
启动成功后, 访问`http://localhost:8559`


#### 方式二、Docker部署【推荐👍】

- 本地配置：以mac os为例
```shell
mkdir -p ~/devman/config
```
将项目中的`./config`目录下的2个配置文件复制过去，并修改其中的数据库配置。修改内容参考上面的方式一

**支持多环境配置**

修改 `boot.yml` 中的`env: dev`示例：
```shell
env:dev  # 对应的config文件为：config-dev.yml
或
env:test  # 对应的config文件为：config-test.yml  
或
env:prod  # 对应的config文件为：config-prod.yml
```


- 拉取镜像
```shell
docker pull jszls65/devman
```
- 运行
```shell
# 将命令中的“~/devman/config”改成你的路径
docker run --name devman_7.1 -p 8559:8559 -v ~/devman/config:/app/config  -d devman
```



#### 方式三、原生部署
1. 本地打包
```shell
# 本地打包，打包成linux上可运行的二进制包
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
```
2. 将devman二进制文件, www和conf两个目录打包上传服务器。

3. 在服务器上解压并运行

服务器上的路径结构
```shell
drwxr-xr-x 2 root root     4096 Jan 16 15:01 config/
-rwxr-xr-x 1 root root 26531274 Jan 16 15:22 devman
drwxr-xr-x 4 root root     4096 Sep 15 10:53 www/
```
启动服务
```shell
nohup ./devman &
```
启动成功后, 访问`http://localhost:8559`


注意：服务器上的config目录应该手动维护，不建议每次打包的时候都将本地的改动覆盖服务器上的配置。所以第一次部署时，请关注服务器上的config目录下配置文件是存在，如果不存在，请手动拷贝项目中的文件，并根据实际情况修改对应的配置。


# 最后
目前仅支持MySQL, 后续会添加其他主流的数据库类型.

欢迎各位大佬提宝贵意见, 也欢迎大家发起 `Merge Request`