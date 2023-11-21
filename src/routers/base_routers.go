package routers

import (
	"github.com/gin-gonic/gin"
)

func BaseRoutersInit(g *gin.Engine) {
	IndexRouterInit(g)
	RequestLogRouterInit(g)
	ToolsRouterInit(g)
	AlertRoutersInit(g)
	NacosRoutersInit(g)
	DatamapRouterInit(g)
}
