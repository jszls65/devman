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
	"time"
)

// 全局变量  记录URL上一次失败时间 map
var lastFailTimeMap = make(map[string]time.Time)

// 服务健康检查
type ServiceAliveCheck struct {
}

func (sa *ServiceAliveCheck) Run() {
	db := persistence.DB
	interfaceConfigList := []*models.InterfaceConfig{}
	db.Where("type = ?", "alive").Find(&interfaceConfigList)
	if len(interfaceConfigList) == 0 {
		log.Println("无接口配置数据")
		return
	}
	for _, interfaceConfig := range interfaceConfigList {
		go sa.doJobsItem(interfaceConfig, db)
		time.Sleep(200 * time.Millisecond)
	}

}

func (sa *ServiceAliveCheck) doJobsItem(interfaceConfig *models.InterfaceConfig, db *gorm.DB) {
	startTime := time.Now()
	httpResult, _ := sa.httpClient(interfaceConfig.HTTPMethod, interfaceConfig.URL)
	log.Println("请求接口是否成功: ", httpResult)
	endTime := time.Now()
	// 计算耗时
	apiDuration := endTime.Sub(startTime)

	lastFailTime, exists := lastFailTimeMap[interfaceConfig.AppName]
	if httpResult == constants.ResultFail { // 请求失败
		interfaceConfig.FailNum += 1 // 失败次数+1
		if !exists {
			lastFailTimeMap[interfaceConfig.AppName] = time.Now()
		} else {
			// 计算失败时间差
			failDuration := time.Now().Sub(lastFailTime)
			if failDuration.Minutes() >= config.Conf.DingTalk.AlertDuration {

				dingTalkMsg := fmt.Sprintf("服务告警: [%s]服务没有响应, 请检查 ! ", interfaceConfig.AppName)
				log.Println("发送钉钉消息: ", dingTalkMsg)
				dingtalk.SendText(dingTalkMsg)
				lastFailTimeMap[interfaceConfig.AppName] = time.Now().Add(config.Conf.DingTalk.NextDuration * time.Minute) // 出现告警后, 10分钟再检查.
			}
		}
	} else { // 请求成功
		if exists {
			delete(lastFailTimeMap, interfaceConfig.AppName)
		}
	}
	interfaceConfig.CallNum += 1            // 总次数+1
	interfaceConfig.UpdateTime = time.Now() // 记录更新时间
	// 更新配置表数据
	db.Where("id = ?", interfaceConfig.ID).Updates(models.InterfaceConfig{CallNum: interfaceConfig.CallNum, UpdateTime: interfaceConfig.UpdateTime})
	// 保存日志
	logItem := &models.InterfaceCallLog{InterfaceConfigID: interfaceConfig.ID, Result: httpResult, CostTime: int32(apiDuration.Milliseconds()), CreateTime: time.Now(), UpdateTime: time.Now()}
	db.Create(logItem)
}

// 发送http请求
func (sa *ServiceAliveCheck) httpClient(method string, url string) (int32, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return constants.ResultFail, err
	}
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Authorization", "Bearer ")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return constants.ResultFail, err
	}
	defer resp.Body.Close()
	return constants.ResultOK, nil
}
