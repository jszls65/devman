package routers

import (
	"devman/src/controllers"

	"github.com/gin-gonic/gin"
)

func DatamapRouterInit(g *gin.Engine) {
	group := g.Group("/datamap")

	group.GET("/", controllers.DatamapController{}.Html)
	group.GET("/refreshCache", controllers.DatamapController{}.RefreshCache)
	group.GET("/load-code", controllers.DatamapController{}.LoadCode)
	group.GET("/table-search", controllers.DatamapController{}.TableSearch)
	group.GET("/share", controllers.DatamapController{}.Share)

}
