package persistence

import (
	"log"

	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"

	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func init() {
	// github.com/mattn/go-sqlite3
	DB, err = gorm.Open(sqlite.Open("./dev-utils.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("连接数据库失败: ", err)
	}
}
