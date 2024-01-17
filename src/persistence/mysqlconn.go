package persistence

import (
	"devman/config"
	"errors"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

// 数据库连接池
var _dbMap map[string]*gorm.DB

func init() {
	mysqlconfigs := config.ListEnableMysqlConfig()
	_dbMap = make(map[string]*gorm.DB)

	for _, mysqlconfig := range mysqlconfigs {
		dsn := mysqlconfig.User + ":" + mysqlconfig.Password + "@tcp(" + mysqlconfig.Host + ":3306)/" + mysqlconfig.DB + "?charset=utf8mb4&parseTime=True&loc=Local"
		_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalln("连接数据库失败: ", err)
		}
		_db.Logger = logger.Default.LogMode(logger.Info)
		sqlDB, _ := _db.DB()
		sqlDB.SetMaxOpenConns(mysqlconfig.MaxOpenConns) //设置数据库连接池最大连接数
		sqlDB.SetMaxIdleConns(mysqlconfig.MaxIdleConns) //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
		_dbMap[mysqlconfig.Env] = _db
	}
}

// 返回特定环境的数据库连接
func GetMysql(name string) (*gorm.DB, error) {
	db, ok := _dbMap[name]
	if !ok {
		return nil, errors.New("env不存在")
	}
	return db, nil
}
