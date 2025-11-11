package config

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Mysql        *mySQL
	Server       *server
	AiEndpoint   *aiEndpoint
	runtimeViper = viper.New()
)

// Init 目的是初始化并读入配置
func Init() {
	configPath := "./config/config.yaml"

	runtimeViper.SetConfigFile(configPath)
	runtimeViper.SetConfigType("yaml")

	if err := runtimeViper.ReadInConfig(); err != nil {
		logger.Fatal("config.Init: 未找到配置文件")
	}

	configMapping()

	// 设置持续监听
	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		// 我们无法确定监听到配置变更时是否已经初始化完毕，所以此处需要做一个判断
		logger.Infof("config: notice config changed: %v\n", e.String())
		configMapping() // 重新映射配置
	})
	runtimeViper.WatchConfig()
}

// configMapping 用于将配置映射到全局变量
func configMapping() {
	c := new(config)
	if err := runtimeViper.Unmarshal(&c); err != nil {
		// 由于这个函数会在配置重载时被再次触发，所以需要判断日志记录方式
		logger.Fatalf("config.configMapping: 配置反序列化失败: %v", err)
	}
	Mysql = &c.MySQL
	Server = &c.Server
	AiEndpoint = &c.AiEndpoint

}
