package controllers

import (
	"devman/src/common/config"
	"devman/src/common/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

type NacosController struct{}

var bootstrapMap = make(map[int]NaocsConfigItem)

// 获取nacos连接
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
func (ic NacosController) Html2GetConfig(c *gin.Context) {
	namespace, _ := c.GetQuery("namespace")
	projectIdOrg, _ := c.GetQuery("proid")
	nacosClient, _ := ic.getNacosClient(namespace)
	projectId, _ := strconv.Atoi(projectIdOrg)

	nacosConfigParams := ic.listNacosConfigParam(projectId)

	var voList = []NaocsConfigItem{}

	for _, val := range nacosConfigParams {
		dataId := val.DataId
		if !strings.Contains(dataId, ".") {
			dataId = dataId + ".properties"
		}

		content, err := nacosClient.GetConfig(vo.ConfigParam{
			DataId: dataId,    //配置文件名
			Group:  val.Group, //配置文件分组
		})
		if err != nil {
			//读取配置文件失败
			log.Println(err)
		}

		// _nacosConfigMap[ic.getKey(val.Group, val.DataId)] = content

		// 文件类型
		fileType := ic.getFileType(dataId)

		voList = append(voList, NaocsConfigItem{
			Name:     ic.getKey(val.Group, dataId),
			FileType: fileType,
			DataId:   dataId,
			Content:  content,
		})
	}

	voListN := make([]NaocsConfigItem, 0)
	voListN = append(voListN, bootstrapMap[projectId])
	voListN = append(voListN, voList...)

	c.HTML(200, "nacos/nacos_config.html", gin.H{
		"nacosConfigs": voListN,
	})
}

func (ic NacosController) getKey(group string, dataId string) string {
	return group + " > " + dataId
}

func (ic NacosController) getFileType(dataId string) string {
	fileType := filepath.Ext(dataId)
	if strings.ToLower(fileType) == ".yml" {
		return "yaml"
	} else if strings.ToLower(fileType) == ".properties" {
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
	DataId   string
	FileType string // 文件类型，前端高亮代码使用
}

// gitlab 项目
type GitLabProject struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sort        int
}

type GitLabProjects []GitLabProject

func (a GitLabProjects) Len() int {
	return len(a)
}

func (a GitLabProjects) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a GitLabProjects) Less(i, j int) bool {
	return a[i].Sort > a[j].Sort // 注意这里是">"，因为我们希望是降序排序
}

// nacos配置参数
type NacosConfigParamBo struct {
	Group  string
	DataId string
}

func (ic NacosController) listNacosConfigParam(projectId int) []NacosConfigParamBo {
	// 根据项目id, 获取项目中bootstrop配置文件内容, 并解析内容得到group 与 dataid信息
	gitProjectInfo := config.GetGitProjectById(projectId)
	if gitProjectInfo.Id == 0 {
		log.Println("id不存在")
		return nil
	}
	bootstrapFilePath := config.Conf.GitLab.Url + "/api/v4/projects/" + strconv.Itoa(projectId) + "/repository/files/" + url.QueryEscape(gitProjectInfo.BootstrapPath) + "?ref=" + gitProjectInfo.Branch
	log.Println("urlStr:", bootstrapFilePath)
	respBodyStr, err := utils.SendHttpRequstGet(bootstrapFilePath, ic.getGitLabHeadMap())
	if err != nil {
		log.Println("获取项目文件内容失败[bootstrap.yaml], 错误内容:", err.Error())
		return nil
	}

	bootstrapFileContent := ic.getFileContent(respBodyStr)

	// 将字符串转换为Reader
	configReader := strings.NewReader(bootstrapFileContent)

	// 初始化Viper
	bootstrapFileType := "properties"
	if strings.Contains(gitProjectInfo.BootstrapPath, ".yml") {
		bootstrapFileType = "yaml"
	}

	bootstrapMap[projectId] = NaocsConfigItem{Name: "bootstrap" + "." + bootstrapFileType,
		Content: bootstrapFileContent, DataId: "bootstrap" + "." + bootstrapFileType, FileType: bootstrapFileType}
	viper.SetConfigType(bootstrapFileType) // 设置配置格式为YAML

	// 从Reader中读取配置
	err = viper.ReadConfig(configReader)
	if err != nil {
		log.Println("viper解析配置变量失败:", err.Error())
		return nil
	}

	// 返回值
	nacosConfigParams := make([]NacosConfigParamBo, 0)

	dataId := viper.GetString("spring.cloud.nacos.config.name")
	if dataId == "${spring.application.name}" {
		dataId = viper.GetString("spring.application.name")
	}
	group := viper.GetString("spring.cloud.nacos.config.group")
	if dataId != "" && group != "" {
		nacosConfigParams = append(nacosConfigParams, NacosConfigParamBo{DataId: dataId, Group: group})
	}

	if bootstrapFileType == "properties" {

		nacosConfigParams = ic.handleNacosKey2Map4Properties("spring.cloud.nacos.config.shared-configs", nacosConfigParams)
		nacosConfigParams = ic.handleNacosKey2Map4Properties("spring.cloud.nacos.config.extension-configs", nacosConfigParams)

	} else {
		nacosConfigParams = ic.handleNacosKey2Map4Yaml("spring.cloud.nacos.config.shared-configs", nacosConfigParams)
		nacosConfigParams = ic.handleNacosKey2Map4Yaml("spring.cloud.nacos.config.extension-configs", nacosConfigParams)
	}

	return nacosConfigParams

}

func (ic NacosController) handleNacosKey2Map4Properties(key string, nacosConfigParams []NacosConfigParamBo) []NacosConfigParamBo {
	for i := 0; i < 10; i++ {
		nacosKey := key + "[" + strconv.Itoa(i) + "]"
		if !viper.IsSet(nacosKey) {
			continue
		}
		server := viper.Get(nacosKey)
		if server == nil {
			break
		}
		serverMap := server.(map[string]interface{})
		dataId := serverMap["data-id"].(string)
		group := serverMap["group"].(string)
		if dataId != "" && group != "" {
			nacosConfigParams = append(nacosConfigParams, NacosConfigParamBo{DataId: dataId, Group: group})
		}
	}
	return nacosConfigParams
}

func (ic NacosController) handleNacosKey2Map4Yaml(key string, paramsList []NacosConfigParamBo) []NacosConfigParamBo {
	if !viper.IsSet(key) {
		return paramsList
	}
	servers := viper.Get(key).([]interface{})
	if len(servers) == 0 {
		return paramsList
	}

	for _, server := range servers {
		serverMap := server.(map[string]interface{})
		dataId := serverMap["data-id"].(string)
		group := serverMap["group"].(string)
		if dataId != "" && group != "" {
			paramsList = append(paramsList, NacosConfigParamBo{DataId: dataId, Group: group})
		}
	}
	return paramsList
}

func (ic NacosController) getGitLabHeadMap() map[string]string {
	headMap := make(map[string]string)
	headMap["PRIVATE-TOKEN"] = config.Conf.GitLab.Token
	return headMap
}

// 解析gitlab返回的文件内容, base64
func (ic NacosController) getFileContent(orgJson string) string {
	var fileInfo GitLabFileInfo
	err := json.Unmarshal([]byte(orgJson), &fileInfo)
	if err != nil || fileInfo.Content == "" {
		log.Println(err.Error())
		return ""
	}

	// 解码Base64字符串
	decodedBytes, err := base64.StdEncoding.DecodeString(fileInfo.Content)
	if err != nil {
		fmt.Println("Error decoding Base64 string:", err)
		return ""
	}
	// 将解码后的字节切片转换为字符串
	return string(decodedBytes)

}

type GitLabFileInfo struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}
