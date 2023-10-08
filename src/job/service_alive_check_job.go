package job

import (
	"dev-utils/config"
	"dev-utils/src/common/constants"
	"dev-utils/src/common/dingtalk"
	"dev-utils/src/persistence"
	"dev-utils/src/persistence/models"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

// 全局变量  记录URL上一次失败时间 map
var lastFailTimeMap = make(map[string]time.Time)

// 服务健康检查
type ServiceAliveCheck struct {
}

func (sa *ServiceAliveCheck) Run() {
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

func (sa *ServiceAliveCheck) doJobsItem(alertJob *models.AlertJob, db *gorm.DB) {
	startTime := time.Now()
	var httpResult int32
	if alertJob.HTTPMethod == "GET" {
		httpResult = sa.httpGet(alertJob.HTTPMethod, alertJob.URL)
	} else {
		httpResult = sa.httpPost(alertJob.HTTPMethod, alertJob.URL, alertJob.Head, alertJob.Body)
	}
	log.Println("请求接口是否成功: ", httpResult)
	endTime := time.Now()
	// 计算耗时
	apiDuration := endTime.Sub(startTime)

	lastFailTime, exists := lastFailTimeMap[alertJob.AppName]
	if httpResult == constants.ResultFail { // 请求失败
		alertJob.FailNum += 1 // 失败次数+1
		if !exists {
			lastFailTimeMap[alertJob.AppName] = time.Now()
		} else {
			// 计算失败时间差
			failDuration := time.Now().Sub(lastFailTime)
			if failDuration.Minutes() >= config.Conf.DingTalk.AlertDuration {

				dingTalkMsg := fmt.Sprintf("服务告警: [%s]服务没有响应, 请检查 ! ", alertJob.AppName)

				dingtalk.SendText(dingTalkMsg)
				lastFailTimeMap[alertJob.AppName] = time.Now().Add(config.Conf.DingTalk.NextDuration * time.Minute) // 出现告警后, 10分钟再检查.
			}
		}
	} else { // 请求成功
		if exists {
			delete(lastFailTimeMap, alertJob.AppName)
		}
	}
	alertJob.CallNum += 1            // 总次数+1
	alertJob.UpdateTime = time.Now() // 记录更新时间
	// 更新配置表数据
	db.Where("id = ?", alertJob.ID).Updates(models.AlertJob{CallNum: alertJob.CallNum, UpdateTime: alertJob.UpdateTime,
		FailNum: alertJob.FailNum})
	// 保存日志
	logItem := &models.AlertLog{AlertJobsID: alertJob.ID, Result: httpResult, CostTime: int32(apiDuration.Milliseconds()), CreateTime: time.Now(), UpdateTime: time.Now()}
	db.Create(logItem)
}

func (sa *ServiceAliveCheck) httpGet(method string, url string) int32 {
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

func (sa *ServiceAliveCheck) httpPost(method string, url string, head string, body string) int32 {
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
