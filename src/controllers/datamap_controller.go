package controllers

import (
	"database/sql"
	"devman/config"
	"devman/src/common"
	"devman/src/persistence"
	"devman/src/structs"
	structsm "devman/src/structs/datamap"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type DatamapController struct{}

var tableInfoMap = make(map[string][]structsm.TableInfo) // 环境名称:数据
var sw = sync.WaitGroup{}                                // 同步等待组
var ch = make(chan int, 1)                               // 协程数量限制
var createTableSqlMap = make(map[string]string)          // 建表语句缓存, key:环境-表名, value是create语句

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

// 刷新缓存
func (ic DatamapController) refreshCache(env string) {
	// 查询数据
	tableInfos := ic.ListTableInfo(env)
	fillTableColumnInfo(env, tableInfos)
	tableInfoMap[env] = tableInfos
}

// 填充表字段信息
func fillTableColumnInfo(env string, infos []structsm.TableInfo) []structsm.TableInfo {
	if len(infos) == 0 {
		return infos
	}
	mysql := persistence.GetMysql(env)
	mysqlByEnv := config.GetMysqlByEnv(env)
	for index, tableItem := range infos {
		tableItem.DbName = mysqlByEnv.DB
		sw.Add(1) // 同步等待组数量
		ch <- 1
		go func(index int, tableItem structsm.TableInfo) {
			defer sw.Done()
			<-ch
			// 表字段切片
			var cols []structsm.ColumnInfo
			// sqlStr := "SELECT \n    COLUMN_NAME Field, IS_NULLABLE `Null`, column_type  `Type`, COLUMN_COMMENT Comment, EXTRA Extra, COLUMN_KEY `Key`, column_default `Default`\nFROM\n    INFORMATION_SCHEMA.COLUMNS\nwhere " +
			// 	"\n    TABLE_SCHEMA = '" + tableItem.DbName + "' and TABLE_NAME = '" + tableItem.TableName + "'\nORDER BY ORDINAL_POSITION;"

			sqlStr := `
				SELECT 
    COLUMN_NAME  TField, IS_NULLABLE TNull, column_type  TType, COLUMN_COMMENT TComment, EXTRA TExtra, COLUMN_KEY TKey, column_default TDefault
    FROM
        INFORMATION_SCHEMA.COLUMNS
        where    TABLE_SCHEMA = @dbName and TABLE_NAME = @tableName
        ORDER BY ORDINAL_POSITION;
				`
			//sqlStr := "desc `" + tableItem.TableName + "`"
			mysql.Raw(sqlStr, sql.Named("dbName", tableItem.DbName), sql.Named("tableName", tableItem.TableName)).Scan(&cols)
			infos[index].Columns = cols
		}(index, tableItem)
	}
	sw.Wait()
	return infos
}

// 查询所有表的表名和表注释列表
func (ic DatamapController) ListTableInfo(env string) []structsm.TableInfo {
	sqlStr := `select
			table_name,
			table_comment
		from information_schema.tables where TABLE_SCHEMA = @dbName;`
	mysql := persistence.GetMysql(env)
	dbName := config.GetMysqlByEnv(env).DB
	// 查询结果
	var tableInfos []structsm.TableInfo
	mysql.Raw(sqlStr, sql.Named("dbName", dbName)).Scan(&tableInfos)
	return tableInfos
}

// 刷新缓存
func (ic DatamapController) RefreshCache(context *gin.Context) {
	env, ok := context.GetQuery("env")
	if !ok {
		env = "测试环境"
	}
	ic.refreshCache(env)
	context.JSON(http.StatusOK, common.ResultMsg(http.StatusOK, "刷新缓存成功"))
}

// 生成建表语句
func (ic DatamapController) LoadCode(context *gin.Context) {
	tableName, _ := context.GetQuery("tableName")
	env, _ := context.GetQuery("env")

	// 获取建表语句
	mysql := persistence.GetMysql(env)
	val, ok := createTableSqlMap[env+"-"+tableName]
	if ok {
		context.HTML(http.StatusOK, "datamap/gencode.html", gin.H{
			"createSql": val,
		})
		return
	}

	// 查询结果
	ret := new(structs.CreateTableBo)
	err := mysql.Raw("show create table " + tableName).Scan(&ret).Error
	if err != nil {
		log.Println("sql执行异常:", err.Error())
	}
	createTableSqlMap[env+"-"+tableName] = ret.CreateTable
	context.HTML(http.StatusOK, "datamap/gencode.html", gin.H{
		"createSql": ret.CreateTable,
	})
}

func (th DatamapController) TableSearch(context *gin.Context) {
	// 获取env
	// 从缓冲中获取表名列表
	env, exists := context.GetQuery("env")
	if !exists {
		panic("参数异常")
	}

	tableInfos := th.ListTableInfo(env)
	tableNames := make([]string, 0)
	for _, info := range tableInfos {
		tableNames = append(tableNames, info.TableName)
	}
	context.HTML(http.StatusOK, "datamap/tablesearch.html", gin.H{
		"tableNames": tableNames,
	})
}

type Result struct {
	ConfigId   int
	ConfigName string
	ConfigKey  string
}
