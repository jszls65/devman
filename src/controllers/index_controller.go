package controllers

import "github.com/gin-gonic/gin"

type IndexController struct{}

func (ic IndexController) Index(c *gin.Context){
	c.HTML(200, "index/index.html", nil);
}