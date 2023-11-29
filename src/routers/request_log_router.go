// @Title
// @Author  zls  2023/9/15 14:05
package routers

import (
	"devman/src/controllers"
	"devman/src/middlewares"

	"github.com/gin-gonic/gin"
)

func RequestLogRouterInit(g *gin.Engine) {
	group := g.Group("/log")
	// 限流中间件
	group.Use(middlewares.RateLimiterMiddleware)

	group.GET("/save", controllers.RequestLogController{}.Save)
	group.GET("/sum", controllers.RequestLogController{}.Sum)
}
