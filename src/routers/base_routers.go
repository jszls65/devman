package routers

import (
	"github.com/gin-gonic/gin"
)

func BaseRoutersInit(g *gin.Engine) {
	AdminRouterInit(g)
	DatamapRouterInit(g)
}
