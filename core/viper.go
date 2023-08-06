package core

import (
	"fmt"
	"github.com/prclin/alumni-circle/config"
	"github.com/prclin/alumni-circle/global"

	"github.com/spf13/viper"
)

/*
loadConfiguration 加载配置文件到Configuration结构体中
*/
func loadConfiguration() {
	v := viper.New()
	v.AddConfigPath("./")
	//此路径为了测试
	v.AddConfigPath("../")
	v.SetConfigName("application")
	v.SetConfigType("yaml")

	//读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	//映射到结构体
	if err := v.Unmarshal(config.DefaultConfiguration); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	global.Configuration = config.DefaultConfiguration
}
