// @Title
// @Author  zls  2023/9/23 16:27
package structs

type PageInfo struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

func GetOffset(page int, limit int) int {
	return (page - 1) * limit
}
