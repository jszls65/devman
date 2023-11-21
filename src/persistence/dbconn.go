package persistence

import (
	"dev-utils/config"
	"gorm.io/gorm/logger"
	"log"

	//"gorm.io/driver/sqlite"  官方的使用cgo, 有时候服务器上并没有cgo环境, 所以改成第三方库了.
	"github.com/glebarez/sqlite"

	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func init() {
	path := config.Conf.SqliteConfig.Path
	log.Println("数据库path:", path)
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalln("连接数据库失败: ", err)
	}
	DB.Logger = logger.Default.LogMode(logger.Error)
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(100) //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭

}
