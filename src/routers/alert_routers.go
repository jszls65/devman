// @Title
// @Author  zls  2023/9/23 23:22
package routers

import (
	"devman/src/controllers"

	"github.com/gin-gonic/gin"
)

func AlertRoutersInit(g *gin.Engine) {
	group := g.Group("/alert")

	group.GET("/load-add", controllers.AlertController{}.LoadAdd)
	group.GET("/load-list", controllers.AlertController{}.AlertListHtml)
	group.GET("/load-edit", controllers.AlertController{}.LoadAdd)
	group.GET("/list", controllers.AlertController{}.GetAlertList)
	group.POST("/add", controllers.AlertController{}.Add)
	group.POST("/del", controllers.AlertController{}.Del)
	group.GET("/update-state", controllers.AlertController{}.UpdateState)
	group.GET("/sdb-open", controllers.AlertController{}.DbOpen)  // sqllite数据库是否开启

}
