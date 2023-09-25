// @Title
// @Author  zls  2023/9/23 23:24
package controllers

import (
	"dev-utils/src/common"
	"dev-utils/src/persistence"
	"dev-utils/src/persistence/models"
	"dev-utils/src/structs"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AlertController struct{}

// 加载html页面
func (ac AlertController) LoadAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "alert/add.html", nil)
}

// 页面
func (ac AlertController) AlertListHtml(c *gin.Context) {
	c.HTML(200, "alert/list.html", gin.H{
		"env": c.Query("dev"),
	})
}

// GetAlertList 查询告警列表
func (ic AlertController) GetAlertList(c *gin.Context) {
	req := &structs.AlertQueryReq{}
	//var req structs.AlertQueryReq
	err := c.BindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ResultMsg("参数不正确"))
		return
	}

	var interfaceConfigs []models.AlertJob
	db := persistence.DB.Model(&interfaceConfigs)
	if req.AppName != "" {
		db.Where("app_name like ?", "%"+req.AppName+"%")
	}
	if req.Owner != "" {
		db.Where("owner like ?", "%"+req.Owner+"%")
	}

	db.Limit(req.Limit).Offset(structs.GetOffset(req.Page, req.Limit))
	db.Find(&interfaceConfigs)

	dataMap := handleDateFormat(interfaceConfigs)
	//c.JSON(200, &interfaceConfigs)
	c.JSON(200, gin.H{
		"code":  0,
		"data":  dataMap,
		"count": len(interfaceConfigs),
	})

}

func (ac AlertController) Add(c *gin.Context) {
	req := new(structs.AlertCreateReq)
	c.ShouldBind(&req)

	job := &models.AlertJob{AppName: req.AppName, HTTPMethod: req.HttpMethod, URL: req.Url, State: req.State,
		Owner: req.Owner, CreateTime: time.Now(), UpdateTime: time.Now()}
	persistence.DB.Create(&job)
	c.JSON(http.StatusOK, common.ResultMsg("操作成功"))

}

func (ac AlertController) Del(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResultMsg("参数异常"))
		return
	}
	paramMap := make(map[string][]int, 0)
	json.Unmarshal(data, &paramMap)

	jobs := []models.AlertJob{}
	persistence.DB.Delete(&jobs, paramMap["ids"])
	c.JSON(http.StatusOK, common.ResultMsg("删除成功"))
}
