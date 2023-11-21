package structs

// 表字段注释 元数据
type TableColumnMatedata struct {
	TableName     string `gorm:"column:TABLE_NAME"`
	ColumnName    string `gorm:"column:COLUMN_NAME"`
	ColumnComment string `gorm:"column:COLUMN_COMMENT"`
}
