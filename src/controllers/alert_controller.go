// @Title
// @Author  zls  2023/9/23 23:24
package controllers

import (
	"dev-utils/src/common"
	"dev-utils/src/common/utils"
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
	idStr, b := c.GetQuery("id")

	job := &models.AlertJob{}
	if b {
		persistence.DB.Where("id=" + idStr).Find(&job)
	}
	c.HTML(http.StatusOK, "alert/add.html", map[string]interface{}{
		"data": job,
	})

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

	var alertJobs []models.AlertJob
	db := persistence.DB.Order("id desc").Model(&alertJobs)
	if req.AppName != "" {
		db.Where("app_name like ?", "%"+req.AppName+"%")
	}
	if req.Owner != "" {
		db.Where("owner like ?", "%"+req.Owner+"%")
	}
	var total int64
	db.Count(&total)

	db.Limit(req.Limit).Offset(structs.GetOffset(req.Page, req.Limit))
	db.Find(&alertJobs)

	voList := handleDateFormat(alertJobs)
	//c.JSON(200, &alertJobs)
	c.JSON(200, gin.H{
		"code":  0,
		"data":  voList,
		"count": total,
	})

}

func (ac AlertController) Add(c *gin.Context) {
	req := new(structs.AlertCreateReq)
	c.ShouldBind(&req)
	// 校验服务名称或url不能重复
	if checkStr := ac.checkAppName(req); checkStr != "" {
		c.JSON(http.StatusInternalServerError, common.ResultMsg(checkStr))
		return
	}
	if checkStr := ac.checkUrl(req); checkStr != "" {
		c.JSON(http.StatusInternalServerError, common.ResultMsg(checkStr))
		return
	}

	job := &models.AlertJob{AppName: req.AppName, HTTPMethod: req.HttpMethod, URL: req.Url, State: req.State,
		Owner: req.Owner, CreateTime: time.Now(), UpdateTime: time.Now(), Note: req.Note, Body: req.Body, Phone: req.Phone}
	var msg string
	if req.Id == 0 {
		// 保存
		persistence.DB.Create(&job)
		msg = "添加成功"
	} else {
		job.ID = req.Id
		persistence.DB.Model(models.AlertJob{}).Where("id", req.Id).Updates(models.AlertJob{AppName: req.AppName,
			State: req.State, HTTPMethod: req.HttpMethod, Owner: req.Owner, URL: req.Url})
		persistence.DB.Exec("UPDATE alert_job SET app_name = ?, owner = ?, state = ?, http_method = ?, url = ?, update_time = ?, note=?, body=? , phone=? WHERE id = ?", req.AppName, req.Owner, req.State, req.HttpMethod, req.Url, time.Now(), req.Note, req.Body, req.Phone, req.Id)
		msg = "编辑成功"
	}
	c.JSON(http.StatusOK, common.ResultMsg(msg))
}

func (ac AlertController) checkUrl(req *structs.AlertCreateReq) string {
	dbJob := &models.AlertJob{}
	sql := persistence.DB.Where("url=?", req.Url) //.Find(&dbJob)
	if req.Id != 0 {
		sql.Where("id != ?", req.Id)
	}
	sql.Find(&dbJob)
	//if dbJob.ID != 0 {
	//	return "url不能重复"
	//}
	// 校验url是否能访问
	errMsg, ok := utils.HttpClient(req.HttpMethod, req.Url, req.Body)
	if !ok {
		return errMsg
	}
	return ""
}

func (ac AlertController) checkAppName(req *structs.AlertCreateReq) string {
	db1 := &models.AlertJob{}
	db1sql := persistence.DB.Where("app_name=?", req.AppName) //.Find(&db1)
	if req.Id != 0 {
		db1sql.Where("id != ?", req.Id)
	}
	db1sql.Find(&db1)
	if db1.ID != 0 {
		return "服务名称不能重复"
	}
	return ""
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

// 更新状态
func (ac AlertController) UpdateState(context *gin.Context) {
	id, exists := context.GetQuery("id")
	if !exists {
		context.JSON(http.StatusBadRequest, common.ResultMsg("id不能为空"))
	}
	state, exists := context.GetQuery("state")
	if !exists {
		context.JSON(http.StatusBadRequest, common.ResultMsg("state不能为空"))
	}
	persistence.DB.Exec("update alert_job set state=? where id=?", state, id)
	context.JSON(http.StatusOK, common.ResultMsg("状态更新成功"))
}
