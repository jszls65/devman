package controllers

import (
	"devman/src/common/constants"
	"devman/src/persistence/models"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type ToolsController struct{}

func (ic ToolsController) Html(c *gin.Context) {
	c.HTML(200, "admin/tools.html", nil)
}

func (ic ToolsController) Welcome(c *gin.Context) {
	c.HTML(200, "admin/welcome.html", nil)
}

// 结构体列表转map,并对日期进行格式化
func handleDateFormat(list []models.AlertJob) []map[string]interface{} {
	mapList := make([]map[string]interface{}, 0)
	for _, val := range list {
		valBytes, _ := json.Marshal(&val)
		dataMap := make(map[string]interface{})
		_ = json.Unmarshal(valBytes, &dataMap)
		dataMap["create_time"] = val.CreateTime.Format(constants.DateTimeMinutes)
		dataMap["update_time"] = val.UpdateTime.Format(constants.DateTimeMinutes)

		mapList = append(mapList, dataMap)
	}

	return mapList
}
