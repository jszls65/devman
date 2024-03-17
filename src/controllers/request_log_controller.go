// @Title
// @Author  zls  2023/9/15 14:06
package controllers

import (
	"devman/config"
	"devman/src/common"
	"devman/src/persistence"
	"devman/src/persistence/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RequestLogController struct {
}

func (s RequestLogController) Save(c *gin.Context) {
	enable := config.Conf.SqliteConfig.Enable
	if !enable {
		// log.Println("sqlite 数据库未启用")
		return
	}
	db := persistence.DB
	// 保存日志数据
	reqLog := models.RequestLog{
		IP:         c.ClientIP(),
		Plateform:  c.GetHeader("User-Agent"),
		Path:       c.FullPath(),
		Event:      c.DefaultQuery("event", ""),
		Params:     "",
		ReqDay:     time.Now().Format(time.DateOnly),
		CreateTime: time.Now(),
	}
	db.Create(&reqLog)
	c.JSON(200, common.ResultOkMsg("操作成功"))
}

// 查询统计数据
func (r RequestLogController) Sum(c *gin.Context) {
	en := config.Conf.SqliteConfig.Enable
	if(!en){
		return
	}
	db := persistence.DB
	var count int64
	db.Model(&models.RequestLog{}).Where("1=1").Count(&count)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": count,
	})
}

var limiter = rate.NewLimiter(rate.Every(time.Second/10), 10)

// 限流器
func (r RequestLogController) limiter() bool {
	return limiter.Allow()
}
