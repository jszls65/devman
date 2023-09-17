package config

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(LibraryConfig)

type MysqlConfig struct {
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

type LibraryConfig struct {
	Mode          string `mapstructure:"mode"`
	Port          int    `mapstructure:"port"`
	*LogConfig    `mapstructure:"log"`
	*MysqlConfig  `mapstructure:"mysql"`
	*RedisConfig  `mapstructure:"redis"`
	*SqliteConfig `mapstructure:"sqlite"`
}

func init() {
	var configPath string

	viper.SetConfigFile("./config/boot.yml")
	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		panic(fmt.Errorf("配置文件读取失败: %s", err))
	}
	configEnv := viper.GetString("env")

	// configEnv := os.Getenv("GO_ENV")
	switch configEnv {
	case "dev":
		configPath = "./config/config-dev.yml"
	case "test":
		configPath = "./config/config-test.yml"
	case "prod":
		configPath = "./config/config-prod.yml"
	default:
		panic("环境变量 GO_ENV 设置错误, 请使用 , 有效值: dev|test|prod, 当前值: " + configEnv)
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
