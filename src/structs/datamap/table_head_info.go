package structs

type TableHeadInfo struct {
	// 数据库名称
	DbName string
	// 表名
	TableName string `gorm:"column:TABLE_NAME"`
	// 表实际备注
	TableComment string `gorm:"column:TABLE_COMMENT"`

	// 数量
	RowCount int

	// 数据容量
	DataVol string

	// 索引容量
	IndexVol string
}
