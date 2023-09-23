// @Title
// @Author  zls  2023/9/23 23:22
package routers

import (
	"dev-utils/src/controllers"
	"github.com/gin-gonic/gin"
)

func AlertRoutersInit(g *gin.Engine) {
	group := g.Group("/alert")

	group.GET("/load-add", controllers.AlertController{}.LoadAdd)
	group.GET("/load-list", controllers.AlertController{}.AlertListHtml)
	group.GET("/list", controllers.AlertController{}.GetAlertList)

}
