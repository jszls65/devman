package job

import (
	"dev-utils/src/common/constants"
	"dev-utils/src/persistence"
	"dev-utils/src/persistence/models"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm/clause"
)

// 服务健康检查

type ServiceAliveCheck struct {
}

func (sa *ServiceAliveCheck) Run() {
	log.Println("执行了")
	db := persistence.DB
	configs := []*models.InterfaceConfig{}
	db.Where("type = ?", "alive").Find(&configs)
	if len(configs) == 0 {
		log.Println("无接口配置数据")
		return
	}
	// 构建日志
	logs := make([]*models.InterfaceCallLog, 0)
	for _, config := range configs {
		startTime := time.Now()
		ok, _ := sa.httpClient(config.HTTPMethod, config.URL)
		endTime := time.Now()
		var re int32 = constants.ResultOK
		if !ok {
			re = constants.ResultFail
			config.FailNum += 1 // 失败次数+1
		}
		config.CallNum += 1            // 总次数+1
		config.UpdateTime = time.Now() // 记录更新时间
		// 计算耗时
		duration := endTime.Sub(startTime)
		logItem := &models.InterfaceCallLog{InterfaceConfigID: config.ID, Result: re, CostTime: int32(duration.Milliseconds()), CreateTime: time.Now(), UpdateTime: time.Now()}

		logs = append(logs, logItem)
	}
	// 更新配置中的调用次数
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"fail_num", "call_num", "update_time"}),
	}).Create(&configs)

	if len(logs) == 0 {
		log.Println("logs无数据")
		return
	}
	// 保存数据
	db.CreateInBatches(logs, 200)

}

func (sa *ServiceAliveCheck) httpClient(method string, url string) (bool, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer ")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}
