// @Title
// @Author  zls  2023/10/12 22:51
package controllers

import (
	"dev-utils/config"
	"dev-utils/src/common"
	"dev-utils/src/structs"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type NacosController struct {
}

var nacos_token_resp structs.NacosTokenResp

// 写个定时任务, 每隔5小时刷新token
func (nc NacosController) Token(context *gin.Context) {
	urlStr := "http://47.97.218.1:8848/nacos/v1/auth/users/login"
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	formData := url.Values{
		"username": []string{"nacos"},
		"password": []string{"asjo123"},
	}
	resp, err := client.PostForm(urlStr, formData)
	if err != nil || resp.StatusCode != 200 {
		log.Println("失败,", err.Error())
	}
	defer resp.Body.Close()
	bodyContent, _ := io.ReadAll(resp.Body)
	log.Println("bodyContent:", string(bodyContent))

	err = json.Unmarshal(bodyContent, &nacos_token_resp)
	if err != nil {
		log.Println("反虚拟化body失败:", err.Error())
		return
	}
	context.JSON(http.StatusOK, common.ResultOk())
}

func (nc NacosController) CheckServiceList(context *gin.Context) {

	healthyServiceMap := getHealthyServiceMap()
	log.Println(healthyServiceMap)

	serviceList := config.Conf.NacosService.List
	log.Println(serviceList)

}

func (nc NacosController) GetToken() string {
	return nacos_token_resp.AccessToken
}

// 获取服务名称与个数的对应关系
func getHealthyServiceMap() map[string]int {
	// http client
	urlStr := "http://47.97.218.1:8848/nacos/v1/ns/catalog/services?hasIpCount=true&withInstances=false&pageNo=1&pageSize=50&accessToken={accessToken}&namespaceId=prod_smart"
	urlStr = strings.ReplaceAll(urlStr, "{accessToken}", nacos_token_resp.AccessToken)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(urlStr)
	if err != nil || resp.StatusCode != 200 {
		log.Println("失败,", err.Error())
	}
	defer resp.Body.Close()
	bodyContent, _ := io.ReadAll(resp.Body)
	// body --> 结构体
	var servicesResp structs.NacosServicesResp
	err = json.Unmarshal(bodyContent, &servicesResp)
	if err != nil {
		log.Println("反虚拟化body失败:", err.Error())
		return nil
	}

	// 对比服务个数是否有缺失
	healthyServicesMap := make(map[string]int) // 健康服务个数
	for _, item := range servicesResp.ServiceList {
		healthyServicesMap[item.Name] = item.HealthyInstanceCount
	}
	return healthyServicesMap

}
