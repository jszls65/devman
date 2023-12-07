// @Title
// @Author  zls  2023/12/6 20:38
package structs

type CreateTableBo struct {
	Table       string `gorm:"column:Table;"`
	CreateTable string `gorm:"column:Create Table;"`
}
