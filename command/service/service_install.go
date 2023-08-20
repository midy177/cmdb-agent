package service

import (
	"cmdb-agent/daemon/service"
	"encoding/base64"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"os"
)

type installCommand struct {
	config service.Config
}

func (c *installCommand) run(*kingpin.ParseContext) error {
	fmt.Printf("installing service %s\n", c.config.Name)
	s, err := service.New(c.config)
	if err != nil {
		return err
	}
	decBytes, err := base64.StdEncoding.DecodeString(c.config.Config)
	if err != nil {
		return err
	}
	err = os.MkdirAll("/etc/cmdb-agent", 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile("/etc/cmdb-agent/config.json", decBytes, 0666)
	if err != nil {
		return err
	}
	return s.Install()
}

func registerInstall(cmd *kingpin.CmdClause) {
	c := new(installCommand)
	s := cmd.Command("install", "install the service").
		Action(c.run)
	s.Flag("config", "client config").
		Required().
		StringVar(&c.config.Config)

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
