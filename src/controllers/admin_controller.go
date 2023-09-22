package controllers

import "github.com/gin-gonic/gin"

type AdminController struct{}

func (ic AdminController) Html(c *gin.Context) {
	c.HTML(200, "admin/admin.html", nil)
}
