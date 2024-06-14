package controllers

import (
	"database/sql"
	"devman/src/common"
	"devman/src/common/config"
	"devman/src/common/utils"
	"devman/src/persistence"
	"devman/src/structs"
	structsm "devman/src/structs/datamap"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type DatamapController struct{}

var tableInfoMap = make(map[string][]structsm.TableInfo) // 环境名称:数据
var sw = sync.WaitGroup{}                                // 同步等待组
var createTableSqlMap = make(map[string]string)          // 建表语句缓存, key:环境-表名, value是create语句

// 主页面
func (ic DatamapController) Html(c *gin.Context) {

	configId, ok := c.GetQuery("configId")
	if !ok {
		c.HTML(200, "datamap/list.html", gin.H{})
		log.Println("参数异常")
		return
	}
	tableInfos, ok := utils.GetMap(tableInfoMap, configId)
	if !ok {
		// 查询数据 刷新缓存
		ic.refreshCache(configId)
		tableInfos, _ = utils.GetMap(tableInfoMap, configId)
	}
	tableNames := make([]string, 0)
	for _, tab := range tableInfos {
		tableNames = append(tableNames, tab.TableName)
	}

	c.HTML(200, "datamap/list.html", gin.H{
		"tableInfos": tableInfos,
		"tableNames": strings.Join(tableNames, ","),
		"configId":   configId,
	})
}

// 刷新缓存
func (ic DatamapController) refreshCache(configId string) {
	// defer func() {
	// 	if re := recover(); re != nil {
	// 		log.Println("刷新缓存失败: ", re)
	// 	}
	// }()
	// 查询数据
	tableInfos := ic.ListTableInfo(configId)
	tableInfos = fillTableColumnInfo(configId, tableInfos)
	//tableInfoMap[env] = tableInfos
	utils.PutMap(tableInfoMap, configId, tableInfos)
}

// 填充表字段信息
func fillTableColumnInfo(configId string, infos []structsm.TableInfo) []structsm.TableInfo {
	if len(infos) == 0 {
		return infos
	}
	mysql, _ := persistence.GetMysql(configId)
	mysqlByEnv := config.GetMysqlByEnv(configId)
	// 对表进行分组
	groupTableInfos := getGroupTableInfos(mysqlByEnv.DB, infos)

	// 返回值
	resultList := []structsm.TableInfo{}
	// 对分组变量进行遍历
	for i := range groupTableInfos {
		tableSubList := groupTableInfos[i]

		sw.Add(1) // 同步等待组数量
		go func(tableItem []structsm.TableInfo) {
			defer sw.Done()
			tableNames := getTableNames(tableSubList)
			// 表字段切片
			var cols []structsm.ColumnInfo
			sqlStr := `
					SELECT 
					TABLE_SCHEMA DbName,  TABLE_NAME TableName,
		COLUMN_NAME  TField, IS_NULLABLE TNull, column_type  TType, COLUMN_COMMENT TComment, EXTRA TExtra, COLUMN_KEY TKey, column_default TDefault
		FROM
			INFORMATION_SCHEMA.COLUMNS
			where    TABLE_SCHEMA = @dbName and TABLE_NAME in @tableNames
			ORDER BY ORDINAL_POSITION;
					`
			//sqlStr := "desc `" + tableItem.TableName + "`"
			mysql.Raw(sqlStr, sql.Named("dbName", tableItem[0].DbName), sql.Named("tableNames", tableNames)).Scan(&cols)
			// infos[index].Columns = cols
			// 按表名分组字段
			groupTableNameColMap := getGroupTableColumns(cols)

			for i := range tableSubList {
				ti := tableSubList[i]
				v, ok := groupTableNameColMap[ti.TableName]
				if !ok {
					continue
				}
				ti.Columns = v
				resultList = append(resultList, ti)
			}
		}(tableSubList)
	}

	sw.Wait()
	// 上面用goroutine异步将表的顺序打乱了, 需要重新排序
	return sortTableList(infos, resultList)
}

func getGroupTableInfos(dbName string, infos []structsm.TableInfo) [][]structsm.TableInfo {
	groupTableInfos := [][]structsm.TableInfo{}
	groupItemTableInfos := []structsm.TableInfo{}
	size := 100
	for index := range infos {
		tableItem := infos[index]
		tableItem.TableHeadInfo.DbName = dbName
		groupItemTableInfos = append(groupItemTableInfos, tableItem)
		if (index+1)%size == 0 {
			// 将分组项数据放入分组变量中
			groupTableInfos = append(groupTableInfos, groupItemTableInfos)
			groupItemTableInfos = []structsm.TableInfo{}
		}
	}
	if len(groupItemTableInfos) > 0 {
		groupTableInfos = append(groupTableInfos, groupItemTableInfos)
	}
	return groupTableInfos
}

// 排序, 按照原始顺序进行排序
func sortTableList(originList []structsm.TableInfo, targetList []structsm.TableInfo) []structsm.TableInfo {
	targetMap := map[string]structsm.TableInfo{}
	for i := range targetList {
		targetItem := targetList[i]
		targetMap[targetItem.TableName] = targetItem
	}
	list := []structsm.TableInfo{}
	for i := range originList {
		tableName := originList[i].TableName
		list = append(list, targetMap[tableName])
	}
	return list
}

func getGroupTableColumns(cols []structsm.ColumnInfo) map[string][]structsm.ColumnInfo {
	groupTableNameColMap := map[string][]structsm.ColumnInfo{}
	for _, ci := range cols {
		tableName := ci.TableName
		val, ok := groupTableNameColMap[tableName]
		if ok {
			val = append(val, ci)
			groupTableNameColMap[tableName] = val
			continue
		}
		// 如果不存在
		firstItem := []structsm.ColumnInfo{}
		firstItem = append(firstItem, ci)
		groupTableNameColMap[tableName] = firstItem
	}
	return groupTableNameColMap
}

// 查询所有表的表名和表注释列表
func (ic DatamapController) ListTableInfo(configId string) []structsm.TableInfo {
	dbName := strings.Split(configId, ",")[1]
	sqlStr := `select
		TABLE_NAME,
		TABLE_COMMENT
		from information_schema.tables where TABLE_SCHEMA = @dbName;`
	mysql, _ := persistence.GetMysql(configId)
	// 查询结果
	var tableMiniInfos []structsm.TableMiniInfo
	mysql.Raw(sqlStr, sql.Named("dbName", dbName)).Scan(&tableMiniInfos)

	var tableInfos []structsm.TableInfo
	for _, info := range tableMiniInfos {
		t := structsm.TableInfo{}
		t.TableName = info.TableName
		if info.TableComment != "" {
			t.TableComment = info.TableComment
		} else {
			t.TableComment = "无备注"
		}

		tableInfos = append(tableInfos, t)
	}
	return tableInfos
}

// 刷新缓存
func (ic DatamapController) RefreshCache(context *gin.Context) {
	configId, ok := context.GetQuery("configId")
	if !ok {
		context.JSON(http.StatusBadRequest, common.ResultMsg(http.StatusBadRequest, "参数错误"))
		return
	}
	ic.refreshCache(configId)
	context.JSON(http.StatusOK, common.ResultMsg(http.StatusOK, "刷新缓存成功"))
}

// 生成建表语句
func (ic DatamapController) LoadCode(context *gin.Context) {
	tableName, _ := context.GetQuery("tableName")
	env, _ := context.GetQuery("env")

	// 获取建表语句
	mysql, _ := persistence.GetMysql(env)
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
	// 获取configId
	// 从缓冲中获取表名列表
	configId, exists := context.GetQuery("configId")
	if !exists {
		panic("参数异常")
	}

	tableInfos, ok := tableInfoMap[configId]
	if !ok {
		time.Sleep(1 * time.Second)
	}
	// tableInfos := th.ListTableInfo(configId)
	tableNames := make([]string, 0)
	for _, info := range tableInfos {
		tableNames = append(tableNames, info.TableName)
	}
	context.HTML(http.StatusOK, "datamap/tablesearch.html", gin.H{
		"tableNames": tableNames,
	})
}

func (ic DatamapController) Share(context *gin.Context) {
	env, _ := context.GetQuery("env")
	tableName, _ := context.GetQuery("tableName")

	//tableInfos, ok := tableInfoMap[env]
	tableInfos, ok := utils.GetMap(tableInfoMap, env)
	filterTableInfos := make([]structsm.TableInfo, 0)
	for _, info := range tableInfos {
		if info.TableName == tableName {
			filterTableInfos = append(filterTableInfos, info)
		}
	}
	if !ok {
		// 查询数据 刷新缓存
		go ic.refreshCache(env)
	}
	tableNames := make([]string, 0)
	for _, tab := range tableInfos {
		tableNames = append(tableNames, tab.TableName)
	}
	context.HTML(200, "datamap/list.html", gin.H{
		"tableInfos": filterTableInfos,
		"tableNames": strings.Join(tableNames, ","),
		"env":        env,
		"shareFlag":  true,
	})
}

type Result struct {
	ConfigId   int
	ConfigName string
	ConfigKey  string
}

func getTableNames(tables []structsm.TableInfo) []string {
	tableNames := []string{}
	for _, item := range tables {
		tableNames = append(tableNames, item.TableName)
	}
	return tableNames
}
