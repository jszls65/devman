package controllers

import (
	"devman/src/common/config"
	"devman/src/common/constants"
	"devman/src/persistence/models"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type AdminController struct{}

func (ic AdminController) Html(c *gin.Context) {
	configs := config.ListEnableMysqlConfig()
	envs := make([]string, 0)
	envGroupMap := make(map[string][]config.MysqlConfig, 0)  // 根据env分组
	for _, v := range configs {
		v.Id = v.Env+","+v.DB
		envs = append(envs, v.Env)
		mapv, ok := envGroupMap[v.Env]
		if ok {
			mapv = append(mapv, v)
			envGroupMap[v.Env] = mapv
			continue
		}
		// 第一次放入map
		firstItem := []config.MysqlConfig{}
		firstItem = append(firstItem, v)
		envGroupMap[v.Env] = firstItem
	}
	c.HTML(200, "admin/admin.html", gin.H{
		"envs": envs,
		"envGroupMap": envGroupMap,
	})
}

func (ic AdminController) Welcome(c *gin.Context) {
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
