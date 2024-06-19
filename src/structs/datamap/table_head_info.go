package structs

type TableHeadInfo struct {
	// 数据库名称
	DbName string `json:"dbName"`
	// 表名
	TableName string `gorm:"column:TABLE_NAME" json:"tableName"`
	// 表实际备注
	TableComment string `gorm:"column:TABLE_COMMENT" json:"tableComment"`

	// 数量
	RowCount int `json:"rowCount"`

	// 数据容量
	DataVol string `json:"dataVol"`
	// 索引容量
	IndexVol string `json:"indexVol"`
}
