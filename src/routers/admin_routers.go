package routers

import (
	"dev-utils/src/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRouterInit(g *gin.Engine) {
	group := g.Group("/admin")

	group.GET("/", controllers.AdminController{}.Html)
	group.GET("/welcome", controllers.AdminController{}.Welcome)
	group.GET("/dl", controllers.FileController{}.DownloadFile)

}
