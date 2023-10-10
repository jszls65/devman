
# Docker
## 移除镜像
docker rmi -f alertman:v1
## 删除容器
docker rm -f alertman
## 构建镜像
docker build -t alertman:v1 .
## 创建容器
docker run -d -p 8559:8559 --name alertman alertman:v1
## 启动容器
docker start alertman

docker cp ../dev-utils.db alertman:/app/alertman

## 导出镜像
docker export alertman > alertman.tar
docker import alertman.tar alertman:v1

# 编译
## 本地环境
go build

## linux 环境
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build