package controllers

import (
	"devman/src/common/config"
	"log"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosController struct{}

// 初始化nacos连接
func (ic NacosController) getNacosClient(namespace string) (config_client.IConfigClient, error) {

	nacosAuth := config.Conf.NacosAuths[0]
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: nacosAuth.IpAddr,       //nacos 地址
			Port:   uint64(nacosAuth.Port), //nacos 端口
		},
	}

	clientConfig := &constant.ClientConfig{
		NamespaceId:         namespace, //命名空间 比较重要 拿取刚才创建的命名空间ID
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "c:/tmp/nacos/log",
		CacheDir:            "c:/tmp/nacos/cache",
		LogLevel:            "debug",
		AccessKey:           nacosAuth.AccessKey,
		SecretKey:           nacosAuth.SecretKey,
	}
	return clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  clientConfig,
			ServerConfigs: serverConfig,
		},
	)
}

// 获取配置
func (ic NacosController) GetConfig(c *gin.Context) {

	namespace, _ := c.GetQuery("namespace")
	group, _ := c.GetQuery("group")

	var dataIds = config.GetNacosDataIds(group)

	// nacos配置map， key: group+'_'+dataId
	var _nacosConfigMap = make(map[string]string)

	nacosClient, _ := ic.getNacosClient(namespace)

	var items = []NaocsConfigItem{}

	for _, val := range dataIds {
		content, err := nacosClient.GetConfig(vo.ConfigParam{
			DataId: val,   //配置文件名
			Group:  group, //配置文件分组
		})
		if err != nil {
			//读取配置文件失败
			log.Println(err)
		}

		_nacosConfigMap[ic.getKey(group, val)] = content

		// 文件类型
		fileType := ic.getFileType(val)

		items = append(items, NaocsConfigItem{
			Name:     ic.getKey(group, val),
			FileType: fileType,
			DataId: val,
			Content:  content,
		})
	}

	c.HTML(200, "nacos/nacos_config.html", gin.H{
		"nacosConfigs": items,
	})
}

func (ic NacosController) getKey(group string, dataId string) string {
	return group + " > " + dataId
}

func (ic NacosController) getFileType(dataId string) string {
	fileType := filepath.Ext(dataId)
	if strings.ToLower(fileType) == ".yml"{
		return "yaml"
	}else if strings.ToLower(fileType) == ".properties"{
		return "properties"
	}
	return "js"
}

// 解析nacos_test.yml配置文件 和上边的字段保持一致
type NacosTestConfig struct {
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	Host string `json:"host" yaml:"host" mapstructure:"host"`
	Port int    `json:"port" yaml:"port" mapstructure:"port"`
}

type NaocsConfigItem struct {
	Name     string
	Content  string
	DataId string  
	FileType string // 文件类型，前端高亮代码使用
}
