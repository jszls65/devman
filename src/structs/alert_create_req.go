// @Title
// @Author  zls  2023/9/25 20:33
package structs

type AlertCreateReq struct {
	Id         int32  `form:"id"`
	AppName    string `form:"appName"`
	HttpMethod string `form:"httpMethod"`
	Url        string `form:"url"`
	Owner      string `form:"owner"`
	State      int32  `form:"state"`
}
