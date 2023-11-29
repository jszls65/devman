package routers

import (
	"devman/src/controllers"

	"github.com/gin-gonic/gin"
)

func IndexRouterInit(g *gin.Engine) {

	g.GET("/index", controllers.IndexController{}.Index)
	g.GET("/tool", controllers.IndexController{}.Index)
}
