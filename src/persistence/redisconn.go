package persistence

import (
	"devman/src/common/config"
	structsm "devman/src/structs/datamap"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
)

// redis客户端
var _rdb *redis.Client

func init() {
	_rdb = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Reids.Addr,     // Redis服务器地址
		Password: config.Conf.Reids.Password, // Redis密码，如果没有设置则留空
		DB:       config.Conf.Reids.DbIndex,  // 选择哪个数据库，默认为0
	})
	log.Println("redis初始化成功")
}

// 返回特定环境的数据库连接
func getClient() *redis.Client {
	return _rdb
}

func SaveData2List(configId string, tableInfos []structsm.TableInfo) error {
	if len(tableInfos) == 0 {
		return nil
	}

	redisKey := getRedisKey4Datamap(configId)
	// 将结构体切片转成string切片
	strList := make([]string, 0)
	for _, val := range tableInfos {
		b, err := json.Marshal(val)
		if err != nil {
			return err
		}
		tableInfoStr := string(b)
		strList = append(strList, tableInfoStr)

	}

	err := getClient().LPush(redisKey, strList).Err()
	if err != nil {
		return err
	}
	// key 设置过期时间
	getClient().Expire(redisKey, time.Hour*6)
	log.Println("数据存入redis成功, redis key: ", redisKey)
	return nil
}

func GetDataFromList(configId string) ([]structsm.TableInfo, error) {
	redisKey := getRedisKey4Datamap(configId)
	// map[string]interface{}

	contentList, err := getClient().LRange(redisKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	tableInfos := make([]structsm.TableInfo, 0)
	for _, contentStr := range contentList {
		var tableInfo structsm.TableInfo
		json.Unmarshal([]byte(contentStr), &tableInfo)
		tableInfos = append(tableInfos, tableInfo)
	}
	return tableInfos, nil
}

func getRedisKey4Datamap(flag string) string {
	return "sqlonedoc:" + flag
}
