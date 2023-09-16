package main

import (
	"dev-utils/src/routers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	// html
	g.LoadHTMLGlob("./www/html/**/*")
	// 静态文件
	g.Static("/static", "./www/static")

	// 中间件
	//g.Use(middlewares.RequestLogMiddleware)

	// 注册路由
	routers.BaseRoutersInit(g)

	// 启动
	log.Println("启动成功")
	err := g.Run(":8559")
	if err != nil {
		log.Fatalln("启动失败: ", err)
		return
	}
}
