// @Title job初始化
// @Author  zls  2023/9/20 23:38
package job

import (
	"devman/config"
	"log"

	"github.com/robfig/cron/v3"
)

var cr *cron.Cron

func init() {
	log.Println("定时任务是否开启: ", config.Conf.Job.Enable)
	if !config.Conf.Job.Enable {
		return
	}
	cr = cron.New(cron.WithSeconds())
	_, err := cr.AddJob(config.Conf.Job.AliveCheck, &ServiceAliveCheck{})
	if err != nil {
		log.Println("定时任务执行失败:", err)
		return
	}
	cr.Start()
}
