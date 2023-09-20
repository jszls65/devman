// @Title job初始化
// @Author  zls  2023/9/20 23:38
package job

import (
	"github.com/robfig/cron/v3"
	"log"
)

var cr *cron.Cron

func init() {
	cr = cron.New()
	_, err := cr.AddJob("*/1 * * * *", &ServiceAliveCheck{})
	if err != nil {
		log.Println("定时任务执行失败:", err)
		return
	}
	cr.Start()
}
