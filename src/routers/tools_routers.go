package routers

import (
	"devman/src/controllers"

	"github.com/gin-gonic/gin"
)

func ToolsRouterInit(g *gin.Engine) {
	group := g.Group("/admin")

	group.GET("/", controllers.ToolsController{}.Html)
	group.GET("/welcome", controllers.ToolsController{}.Welcome)
	group.GET("/dl", controllers.FileController{}.DownloadFile)

}
