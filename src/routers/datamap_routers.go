package routers

import (
	"dev-utils/src/controllers"

	"github.com/gin-gonic/gin"
)

func DatamapRouterInit(g *gin.Engine) {
	group := g.Group("/datamap")

	group.GET("/", controllers.DatamapController{}.Html)
	group.GET("/refreshCache", controllers.DatamapController{}.RefreshCache)

}
