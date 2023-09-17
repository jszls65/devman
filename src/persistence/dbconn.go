package persistence

import (
	"dev-utils/config"
	"log"

	//"gorm.io/driver/sqlite"  官方的使用cgo, 有时候服务器上并没有cgo环境, 所以改成第三方库了.
	"github.com/glebarez/sqlite"

	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func init() {
	log.Println("dbconn init")
	path := config.Conf.SqliteConfig.Path
	log.Println("数据库path:", path)
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalln("连接数据库失败: ", err)
	}
}
