package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Config config

// Config stores the system configuration.
type config struct {
	Logger struct {
		Level      string `mapstructure:"level"`
		File       string `mapstructure:"file"`
		MaxAge     int    `mapstructure:"maxAge"`
		MaxBackups int    `mapstructure:"maxBackups"`
		MaxSize    int    `mapstructure:"maxSize"`
	} `mapstructure:"logger"`
}

func FromYamlFile() {
	viper.SetConfigName("config")            // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")              // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("/etc/cmdb-agent/")  // 查找配置文件所在的路径
	viper.AddConfigPath("$HOME/.cmdb-agent") // 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")                 // 还可以在工作目录中查找配置
	viper.SetDefault("logger.level", "error")
	viper.SetDefault("logger.file", "/var/log/cmdb-agent/run.log")
	viper.SetDefault("logger.maxAge", 1)
	viper.SetDefault("logger.maxBackups", 1)
	viper.SetDefault("logger.maxSize", 100)
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
