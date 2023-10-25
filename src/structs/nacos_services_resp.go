// @Title
// @Author  zls  2023/10/13 07:46
package structs

type NacosServicesResp struct {
	Count       int                     `json:"count"`
	ServiceList []NacosServicesItemResp `json:"serviceList"`
}

type NacosServicesItemResp struct {
	Name                 string `json:"name"`
	GroupName            string `json:"groupName"`
	ClusterCount         int    `json:"clusterCount"`
	IpCount              int    `json:"ipCount"`
	HealthyInstanceCount int    `json:"healthyInstanceCount"`
	TriggerFlag          bool   `json:"triggerFlag"`
}
