package service

import (
	"cmdb-agent/daemon/service"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
)

type startCommand struct {
	config service.Config
}

func (c *startCommand) run(*kingpin.ParseContext) error {
	fmt.Printf("starting service %s\n", c.config.Name)
	s, err := service.New(c.config)
	if err != nil {
		return err
	}
	return s.Start()
}

func registerStart(cmd *kingpin.CmdClause) {
	c := new(startCommand)
	s := cmd.Command("start", "start the service").
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
