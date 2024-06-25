// @Title
// @Author  zls  2023/12/6 20:38
package structs

type CreateTableBo struct {
	Table       string `gorm:"column:Table;"`
	View       string `gorm:"column:View;"`
	CreateTable string `gorm:"column:Create Table;"`
	CreateView string `gorm:"column:Create View;"`
}
