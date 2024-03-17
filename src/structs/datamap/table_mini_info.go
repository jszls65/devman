package structs

type TableMiniInfo struct {

	// 表名
	TableName string `gorm:"column:TABLE_NAME"`
	// 表实际备注
	TableComment string `gorm:"column:TABLE_COMMENT"`
}