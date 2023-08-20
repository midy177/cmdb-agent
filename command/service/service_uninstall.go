package service

import (
	"cmdb-agent/daemon/service"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"os"
)

type uninstallCommand struct {
	config service.Config
}

func (c *uninstallCommand) run(*kingpin.ParseContext) error {
	fmt.Printf("uninstalling service %s\n", c.config.Name)
	s, err := service.New(c.config)
	if err != nil {
		return err
	}
	fmt.Println("开始删除日志文件...")
	err0 := os.RemoveAll("/var/log/cmdb-agent")
	if err0 != nil {
		fmt.Printf("删除失败，%v", err0)
	}
	fmt.Println("开始删除配置文件...")
	err1 := os.RemoveAll("/etc/cmdb-agent")
	if err1 != nil {
		fmt.Printf("删除失败，%v", err1)
	}
	return s.Uninstall()
}

func registerUninstall(cmd *kingpin.CmdClause) {
	c := new(uninstallCommand)
	s := cmd.Command("uninstall", "uninstall the service").
		Action(c.run)

	s.Flag("name", "service name").
		Default(DefaultName).
		StringVar(&c.config.Name)

	s.Flag("desc", "service description").
		Default(DefaultDesc).
		StringVar(&c.config.Desc)

	s.Flag("username", "windows account username").
		Default("").
		StringVar(&c.config.Username)

	s.Flag("password", "windows account password").
		Default("").
		StringVar(&c.config.Password)
}
