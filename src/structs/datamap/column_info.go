package structs

type ColumnInfo struct {
	// 列名
	TField string
	// 数据类型
	TType string
	// 是否可为null
	TNull string
	// 主键
	TKey string
	// 默认值
	TDefault string
	// 描述
	TExtra string
	// 备注
	TComment string
}
