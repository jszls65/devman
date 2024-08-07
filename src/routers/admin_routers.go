package routers

import (
	"devman/src/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRouterInit(g *gin.Engine) {
	group := g.Group("/")

	group.GET("/", controllers.AdminController{}.Html)
	group.GET("/welcome", controllers.AdminController{}.Welcome)

	// nacos
	group.GET("/nacos_config", controllers.NacosController{}.Html2GetConfig)
	group.GET("/nacos_discovery", controllers.NacosController{}.Html2Discovery)
	group.GET("/nacos_discovery/list", controllers.NacosController{}.DiscoveryList)
	group.GET("/nacos_discovery/enable", controllers.NacosController{}.DiscoveryEnable)
}
