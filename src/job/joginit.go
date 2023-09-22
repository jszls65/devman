// @Title job初始化
// @Author  zls  2023/9/20 23:38
package job

import (
	"dev-utils/config"
	"github.com/robfig/cron/v3"
	"log"
)

var cr *cron.Cron

func init() {
	log.Println("定时任务是否开启: ", config.Conf.Job.Enable)
	if !config.Conf.Job.Enable {
		return
	}
	cr = cron.New(cron.WithSeconds())
	_, err := cr.AddJob("*/10 * * * * *", &ServiceAliveCheck{})
	if err != nil {
		log.Println("定时任务执行失败:", err)
		return
	}
	cr.Start()
}
