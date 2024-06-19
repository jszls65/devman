package structs

type TableInfo struct {
	TableHeadInfo `json:"tableHeadInfo"`

	// 表列明细信息
	Columns []ColumnInfo `json:"columns"`
	// 表分组
	Group []string `json:"group"`
	// 分表数统计
	SplitCount int32 `json:"splitCount"`
}
