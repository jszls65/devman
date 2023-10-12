// @Title
// @Author  zls  2023/9/16 10:26
package middlewares

import (
	"dev-utils/src/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RateLimiterMiddleware(c *gin.Context) {
	limiter := common.Limiter(c.FullPath(), 10, 30)
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, common.ResultMsg(1, "请求太频繁"))
		return
	}
	c.Next()

}
