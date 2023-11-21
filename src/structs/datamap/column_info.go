package structs

type ColumnInfo struct {
	// 列名
	Field string
	// 数据类型
	Type string
	// 是否可为null
	IsNull string
	// 主键
	Key string
	// 默认值
	Default string
	// 描述
	Extra string
	// 备注
	Comment string
}
