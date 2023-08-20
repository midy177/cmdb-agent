package service

import (
	"cmdb-agent/daemon/service"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
)

type stopCommand struct {
	config service.Config
}

func (c *stopCommand) run(*kingpin.ParseContext) error {
	fmt.Printf("stopping service %s\n", c.config.Name)
	s, err := service.New(c.config)
	if err != nil {
		return err
	}
	return s.Stop()
}

func registerStop(cmd *kingpin.CmdClause) {
	c := new(stopCommand)
	s := cmd.Command("stop", "stop the service").
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
