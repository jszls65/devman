package main

import (
	"devman/src/common/config"
	"devman/src/routers"
	"devman/src/templatefuns"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode) // todo 测试代码
	g := gin.Default()
	
	// 初始化模板函数
	templatefuns.InitTemplateHandler(g)
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
