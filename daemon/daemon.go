package daemon

import (
	"cmdb-agent/client"
	"cmdb-agent/config"
	"context"
	"github.com/orandin/lumberjackrus"
	"github.com/sirupsen/logrus"
	"log"
)

// Run runs the service and blocks until complete.
func Run(ctx context.Context) error {
	config.FromYamlFile()
	setupLogger()
	//return client.Run(ctx)
	return client.NewClient(ctx)
}

// helper function configures the global logger from
// the loaded configuration.
func setupLogger() {
	level := logrus.InfoLevel
	err := level.UnmarshalText([]byte(config.Config.Logger.Level))
	if err != nil {
		log.Fatalf("level值可以是(panic,fatal,error,warn,warning,info,debug,trace),默认是:trace,解析level出错: %s", err)
	}
	logrus.SetLevel(level)
	if config.Config.Logger.File == "" {
		log.Fatalf("日志文件定义错误！！")
	}
	hook, err := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{
			Filename:   config.Config.Logger.File,
			MaxSize:    config.Config.Logger.MaxSize,
			MaxBackups: config.Config.Logger.MaxBackups,
			MaxAge:     config.Config.Logger.MaxAge,
		},
		level,
		&logrus.TextFormatter{},
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}
	logrus.AddHook(hook)
}
