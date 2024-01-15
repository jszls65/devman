package config

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(LibraryConfig)

type MysqlConfig struct {
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

type SqliteConfig struct {
	Path string `mapstructure:"path" json:"path"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	Password     string `mapstructure:"password"`
	PollSize     int    `mapstructure:"PollSize"`
	MinIdleConns int    `mapstructure:"min_idle_cons"`
}
type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackUps int    `mapstructure:"max_backups"`
}

type DingTalk struct {
	Url           string        `mapstructure:"url"`
	AlertDuration float64       `mapstructure:"alert_duration"`
	NextDuration  time.Duration `mapstructure:"next_duration"`
	Enable        bool          `mapstructure:"enable"`
}

type Job struct {
	Enable     bool   `mapstructure:"enable"`
	AliveCheck string `mapstructure:"alive_check"`
	AutoRule   string `mapstructure:"auto_rule"`
}

type NacosService struct {
	List []map[string]int `mapstructure:"list"`
}

// 总结构体
type LibraryConfig struct {
	Mode          string `mapstructure:"mode"`
	Port          int    `mapstructure:"port"`
	*LogConfig    `mapstructure:"log"`
	MysqlConfigs  []MysqlConfig `mapstructure:"mysqls"`
	*RedisConfig  `mapstructure:"redis"`
	*SqliteConfig `mapstructure:"sqlite"`
	*DingTalk     `mapstructure:"ding_talk"`
	*Job          `mapstructure:"job"`
	*NacosService `mapstructure:"nacos_service"`
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
func GetMysqlByEnv(env string) *MysqlConfig {
	for _, conf := range Conf.MysqlConfigs {
		if conf.Env == env && conf.Enable {
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
