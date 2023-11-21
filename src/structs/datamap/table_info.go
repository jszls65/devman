package structs

type TableInfo struct {
	TableHeadInfo

	// 表列明细信息
	Columns []ColumnInfo
	// 表分组
	Group []string
	// 分表数统计
	SplitCount int32
}
