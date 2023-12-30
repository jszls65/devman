package main

import (
	"devman/config"
	_ "devman/src/job"
	"devman/src/routers"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	//func main4() {
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
	err := g.Run(":" + strconv.Itoa(config.Conf.Port))
	if err != nil {
		log.Fatalln("启动失败: ", err)
		return
	}
}
