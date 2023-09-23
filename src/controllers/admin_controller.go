package controllers

import (
	"dev-utils/src/common/constants"
	"dev-utils/src/persistence/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type AdminController struct{}

func (ic AdminController) Html(c *gin.Context) {
	c.HTML(200, "admin/admin.html", nil)
}

func (ic AdminController) Welcome(c *gin.Context) {
	c.HTML(200, "admin/welcome.html", nil)
}

// 结构体列表转map,并对日期进行格式化
func handleDateFormat(list []models.InterfaceConfig) []map[string]interface{} {
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
