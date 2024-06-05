package templatefuns

import (
	"github.com/gin-gonic/gin"
	"html/template"
)

// InitTemplateHandler 初始化模板函数
func InitTemplateHandler(g *gin.Engine) {
	g.SetFuncMap(template.FuncMap{
		"indexAddOne": indexAddOneTemplateHandler,
	})
}

// indexAddOneTemplateHandler 表格字段序号处理
func indexAddOneTemplateHandler(index int) int {
	return index + 1
}
