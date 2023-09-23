package routers

import (
	"github.com/gin-gonic/gin"
)

func BaseRoutersInit(g *gin.Engine) {
	IndexRouterInit(g)
	RequestLogRouterInit(g)
	AdminRouterInit(g)
	AlertRoutersInit(g)
}
