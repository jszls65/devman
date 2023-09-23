// @Title
// @Author  zls  2023/9/23 16:24
package structs

// 告警列表查询
type AlertQueryReq struct {
	PageInfo
	Env     string `json:"env" form:"env"`
	AppName string `json:"appName" form:"appName"`
	Owner   string `json:"owner" form:"owner"`
}
