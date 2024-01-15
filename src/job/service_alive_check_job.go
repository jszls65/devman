package job

import (
	"devman/config"
	"devman/src/common/constants"
	"devman/src/common/dingtalk"
	"devman/src/persistence"
	"devman/src/persistence/models"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

// 服务健康检查
type ServiceAliveCheckJob struct {
}

func (sa *ServiceAliveCheckJob) Run() {
	if !config.Conf.Job.Enable {
		log.Println("job enable=false, 任务停止")
		return
	}

	db := persistence.DB
	var interfaceConfigList []*models.AlertJob
	db.Where("type = ? and state = 1", "alive").Find(&interfaceConfigList)
	if len(interfaceConfigList) == 0 {
		log.Println("无接口配置数据")
		return
	}
	for _, interfaceConfig := range interfaceConfigList {
		go sa.doJobsItem(interfaceConfig, db)
	}

}

func (sa *ServiceAliveCheckJob) doJobsItem(alertJob *models.AlertJob, db *gorm.DB) {
	// startTime := time.Now()
	var httpResult int32
	if alertJob.HTTPMethod == "GET" {
		httpResult = sa.httpGet(alertJob.URL)
	} else {
		httpResult = sa.httpPost(alertJob.URL, alertJob.Body)
	}

	lastHeathState := alertJob.HeathState
	if httpResult == constants.ResultFail { // 请求失败
		if alertJob.LastFailTime.IsZero() || lastHeathState == 0 {
			alertJob.HeathState = 1            // 健康状态 0-健康 1-告警 2-离线
			alertJob.LastFailTime = time.Now() // 记录失败时间
		} else {
			// 计算失败时间差
			failDuration := time.Since(alertJob.LastFailTime)
			if failDuration.Minutes() >= config.Conf.DingTalk.AlertDuration && alertJob.HeathState != 2 {
				dingTalkMsg := fmt.Sprintf("服务告警: [%s]服务没有响应, 请检查 ! ", alertJob.AppName)

				dingtalk.SendText(dingTalkMsg, alertJob.Phone)
				alertJob.HeathState = 2
			}

		}
	} else { // 请求成功
		alertJob.HeathState = 0
	}
	alertJob.UpdateTime = time.Now() // 记录更新时间
	// 更新配置表数据
	db.Save(&alertJob)
	// 保存日志
	// logItem := &models.AlertLog{AlertJobsID: alertJob.ID, Result: httpResult, CostTime: int32(apiDuration.Milliseconds()), CreateTime: time.Now(), UpdateTime: time.Now()}
	// db.Create(logItem)
}

func (sa *ServiceAliveCheckJob) httpGet(url string) int32 {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return constants.ResultFail
	}
	defer resp.Body.Close()
	return constants.ResultOK
}

func (sa *ServiceAliveCheckJob) httpPost(url string, body string) int32 {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	client.Head("")
	resp, err := client.Post(url, "application/json", strings.NewReader(body))
	if err != nil || resp.StatusCode != 200 {
		return constants.ResultFail
	}
	defer resp.Body.Close()
	return constants.ResultOK
}
