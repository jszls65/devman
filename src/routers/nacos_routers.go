// @Title
// @Author  zls  2023/9/23 23:22
package routers

import (
	"devman/src/controllers"

	"github.com/gin-gonic/gin"
)

func NacosRoutersInit(g *gin.Engine) {
	group := g.Group("/nacos")

	group.GET("/token", controllers.NacosController{}.Token)
	group.GET("/check-service-list", controllers.NacosController{}.CheckServiceList)
}
