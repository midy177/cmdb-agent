package service

import (
	"cmdb-agent/daemon/service"
	"github.com/alecthomas/kingpin/v2"
)

type runCommand struct {
	config service.Config
}

func (c *runCommand) run(k *kingpin.ParseContext) error {
	s, err := service.New(c.config)
	if err != nil {
		return err
	}
	return s.Run()
}

func registerRun(cmd *kingpin.CmdClause) {
	c := new(runCommand)
	s := cmd.Command("run", "run the service").
		Action(c.run).Hidden()

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
