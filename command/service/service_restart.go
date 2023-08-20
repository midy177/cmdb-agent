package service

import (
	"cmdb-agent/daemon/service"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
)

type restartCommand struct {
	config service.Config
}

func (c *restartCommand) run(*kingpin.ParseContext) error {
	fmt.Printf("restart the service %s......\n", c.config.Name)
	s, err := service.New(c.config)
	if err != nil {
		return err
	}
	return s.Restart()
}

func registerRestart(cmd *kingpin.CmdClause) {
	c := new(restartCommand)
	s := cmd.Command("restart", "restart the service").
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
