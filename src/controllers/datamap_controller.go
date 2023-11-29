package controllers

import (
	"database/sql"
	"devman/config"
	"devman/src/common"
	"devman/src/persistence"
	structs "devman/src/structs/datamap"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type DatamapController struct{}

var tableInfoMap = make(map[string][]structs.TableInfo) // 环境名称:数据
var sw = sync.WaitGroup{}                               // 同步等待组
var ch = make(chan int, 1)                              // 协程数量限制

func (ic DatamapController) Html(c *gin.Context) {

	env, ok := c.GetQuery("env")
	if !ok {
		env = "生产环境"
	}
	tableInfos, ok := tableInfoMap[env]
	if !ok {
		// 查询数据 刷新缓存
		go ic.refreshCache(env)
	}
	tableNames := make([]string, 0)
	for _, tab := range tableInfos {
		tableNames = append(tableNames, tab.TableName)
	}
	c.HTML(200, "datamap/list.html", gin.H{
		"tableInfos": tableInfos,
		"tableNames": strings.Join(tableNames, ","),
		"env":        env,
	})
}

func (ic DatamapController) refreshCache(env string) {
	// 查询数据
	tableInfos := ic.ListTableInfo(env, config.GetMysqlByEnv(env).DB)
	fillTableColumnInfo(env, tableInfos)
	tableInfoMap[env] = tableInfos
}


func fillTableColumnInfo(env string, infos []structs.TableInfo) []structs.TableInfo {
	if len(infos) == 0 {
		return infos
	}
	mysql := persistence.GetMysql(env)
	mysqlByEnv := config.GetMysqlByEnv(env)
	for index, tableItem := range infos {
		tableItem.DbName = mysqlByEnv.DB
		sw.Add(1) // 同步等待组数量
		ch <- 1
		go func(index int, tableItem structs.TableInfo) {
			defer sw.Done()
			<-ch
			// 表字段切片
			var cols []structs.ColumnInfo
			sqlStr := "SELECT \n    COLUMN_NAME Field, IS_NULLABLE `Null`, column_type  `Type`, COLUMN_COMMENT Comment, EXTRA Extra, COLUMN_KEY `Key`, column_default `Default`\nFROM\n    INFORMATION_SCHEMA.COLUMNS\nwhere " +
				"\n    TABLE_SCHEMA = '" + tableItem.DbName + "' and TABLE_NAME = '" + tableItem.TableName + "'\nORDER BY ORDINAL_POSITION;"
			//sqlStr := "desc `" + tableItem.TableName + "`"
			mysql.Raw(sqlStr).Scan(&cols)
			infos[index].Columns = cols
		}(index, tableItem)
	}
	sw.Wait()
	return infos
}

func (ic DatamapController) ListTableInfo(env string, dbName string) []structs.TableInfo {
	sqlStr := `select
			table_name,
			table_comment
		from information_schema.tables where TABLE_SCHEMA = @dbName;`
	mysql := persistence.GetMysql(env)
	// 查询结果
	var tableInfos []structs.TableInfo
	mysql.Raw(sqlStr, sql.Named("dbName", dbName)).Scan(&tableInfos)
	return tableInfos
}

func (ic DatamapController) RefreshCache(context *gin.Context) {
	env, ok := context.GetQuery("env")
	if !ok {
		env = "测试环境"
	}
	ic.refreshCache(env)
	context.JSON(http.StatusOK, common.ResultMsg(http.StatusOK, "刷新缓存成功"))
}

type Result struct {
	ConfigId   int
	ConfigName string
	ConfigKey  string
}
