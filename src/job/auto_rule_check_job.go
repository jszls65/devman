package job

import (
	"database/sql"
	"devman/src/common/config"
	"devman/src/common/dingtalk"
	"devman/src/persistence"
	"log"
	"strconv"
	"time"
)

// 服务健康检查
type AutoRuleCheckJob struct {
}

func (sa *AutoRuleCheckJob) Run() {
	if !config.Conf.Job.Enable {
		log.Println("job enable=false, 任务停止")
		return
	}
	now := time.Now()
	hour := now.Hour()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
	startTimeStr := startTime.Format(time.DateTime)
	endTimeStr := time.Now().Format(time.DateTime)

	db, err := persistence.GetMysql("生产环境")
	if err != nil {
		log.Println(err.Error())
		return
	}
	execedNum := 0 // 已经执行的记录数
	db.Raw(`select count(*) from t_amz_adv_auto_log 
		where create_time >= @startTime and  create_time <= @endTime`,
		sql.Named("startTime", startTimeStr),
		sql.Named("endTime", endTimeStr)).Scan(&execedNum)
	log.Printf("[%s -- %s] 时间内, 有%d个类型被执行", startTimeStr, endTimeStr, execedNum)
	if execedNum == 0 {
		// 发送钉钉消息
		dingtalk.SendText("服务告警: [自动化服务]没有执行", "15651636203")
		return
	}

	// 检查失败记录数
	failNum := 0
	db.Raw(`select count(*) from t_amz_adv_auto_log 
		where create_time >= @startTime and  create_time <= @endTime
		and result != '执行成功'
		`,
		sql.Named("startTime", startTimeStr),
		sql.Named("endTime", endTimeStr)).Scan(&failNum)
	log.Printf("[%s -- %s] 时间内, 失败记录数: %d", startTimeStr, endTimeStr, failNum)
	if failNum > 0 {
		// 发送钉钉消息
		dingtalk.SendText("服务告警: [自动化服务]失败: "+strconv.Itoa(failNum)+" 条", "15651636203")
	}

}
