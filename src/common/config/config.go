package config

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(LibraryConfig)

// 总结构体
type LibraryConfig struct {
	Port         int              `mapstructure:"port"`
	MysqlConfigs []MysqlConfig    `mapstructure:"mysqls"`
	NacosAuths   []NacosAuth      `mapstructure:"nacos_auths"`
	NacosGroups  []NacosGroupInfo `mapstructure:"nacos_groups"`
	Reids        RedisConfig      `mapstructure:"redis"`
	GitLab       GitLabConfig     `mapstructure:"gitlab"`
}

type GitLabConfig struct {
	Url      string                `mapstructure:"url"`
	Token    string                `mapstructure:"token"`
	Projects []GitLabProjectConfig `mapstructure:"projects"`
}

type GitLabProjectConfig struct {
	Id int `mapstructure:"id"`
	Name string `mapstructure:"name"`
	BootstrapPath string `mapstructure:"bootstrap_path"`
	Branch string `mapstructure:"branch"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DbIndex  int    `mapstructure:"dbIndex"`
}

type MysqlConfig struct {
	Id           string
	Env          string `mapstructure:"env"`
	Enable       bool   `mapstructure:"enable"`
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"db"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type NacosAuth struct {
	Env       string `mapstructure:"env"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	IpAddr    string `mapstructure:"ipAddr"`
	Port      int    `mapstructure:"port"`
}

type NacosGroupInfo struct {
	Group      string   `mapstructure:"group"`
	DataIds    []string `mapstructure:"dataIds"`
	ServerName string   `mapstructure:"serverName"`
}

func init() {
	//var configPath string
	
	viper.SetConfigFile("./config/boot.yml")
	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		panic(fmt.Errorf("配置文件读取失败: %s", err))
	}
	env := viper.GetString("env")
	configPath, ok := getEnvConfigMap(env)
	if !ok {
		panic("环境变量 GO_ENV 设置错误, 请使用 , 有效值: dev|test|prod, 当前值: " + env)
	}
	log.Println("当前加载的配置文件是: ", configPath)

	//加载配置文件位置
	viper.SetConfigFile(configPath)
	//监听配置文件
	viper.WatchConfig()
	//监听是否更改配置文件
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件被人修改了...")
		err := viper.Unmarshal(&Conf)
		if err != nil {
			panic(fmt.Errorf("配置文件修改以后, 报错啦, err:%v", err))
		}
	})
	// 读取配置文件内容
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("ReadInConfig failed,err:%v", err))
	}
	//将配置文件内容写入到Conf结构体
	if err1 := viper.Unmarshal(&Conf); err1 != nil {
		panic(fmt.Errorf("unmarshal data to Conf failed,err:%v", err))
	}
}

var envMap = make(map[string]string)

func getEnvConfigMap(key string) (string, bool) {
	if len(envMap) == 0 {
		envMap["dev"] = "./config/config-dev.yml"
		envMap["test"] = "./config/config-test.yml"
		envMap["prod"] = "./config/config-prod.yml"
	}
	val, ok := envMap[key]
	return val, ok
}

// 根据环境变量名称获取mysql的配置
func GetMysqlByEnv(configId string) *MysqlConfig {
	for _, conf := range Conf.MysqlConfigs {
		conf.Id = conf.Env + "," + conf.DB
		if conf.Id == configId && conf.Enable {

			return &conf
		}
	}
	return nil
}

// 获取有效mysql配置
func ListEnableMysqlConfig() []MysqlConfig {
	list := make([]MysqlConfig, 0)
	for _, conf := range Conf.MysqlConfigs {
		if !conf.Enable {
			continue
		}
		list = append(list, conf)
	}
	return list
}

func GetNacosDataIds(group string) []string {
	for _, v := range Conf.NacosGroups {
		if v.Group == group {
			return v.DataIds
		}
	}
	return make([]string, 0)
}


// 获取配置中的gitlab项目列表
func ListGitProjects() []GitLabProjectConfig{
	return Conf.GitLab.Projects
}

func GetGitProjectById(id int) GitLabProjectConfig{
	for _, pro := range Conf.GitLab.Projects {
		if pro.Id == id{
			return pro
		}
	}
	return GitLabProjectConfig{}
}