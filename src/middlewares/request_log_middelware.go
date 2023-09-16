// @Title
// @Author  zls  2023/9/15 14:24
package middlewares

import (
	"dev-utils/src/controllers"
	"github.com/gin-gonic/gin"
)

func RequestLogMiddleware(c *gin.Context) {
	c.Next()
	controllers.RequestLogController{}.Save(c)
}
