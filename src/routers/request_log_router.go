// @Title
// @Author  zls  2023/9/15 14:05
package routers

import (
	"dev-utils/src/controllers"
	"dev-utils/src/middlewares"
	"github.com/gin-gonic/gin"
)

func RequestLogRouterInit(g *gin.Engine) {
	group := g.Group("/log")
	// 限流中间件
	group.Use(middlewares.RateLimiterMiddleware)

	group.GET("/save", controllers.RequestLogController{}.Save)
	group.GET("/sum", controllers.RequestLogController{}.Sum)
}
