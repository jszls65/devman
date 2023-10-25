// @Title
// @Author  zls  2023/10/13 07:27
package structs

type NacosTokenResp struct {
	AccessToken string `json:"accessToken"`
	GlobalAdmin bool   `json:"globalAdmin"`
	TokenTtl    int    `json:"tokenTtl"`
	Username    string `json:"username"`
}
