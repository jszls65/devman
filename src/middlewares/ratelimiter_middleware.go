// @Title
// @Author  zls  2023/9/16 10:26
package middlewares

import (
	"devman/src/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware(c *gin.Context) {
	limiter := common.Limiter(c.FullPath(), 10, 30)
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, common.ResultMsg(1, "请求太频繁"))
		return
	}
	c.Next()

}
