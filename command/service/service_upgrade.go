package service

import (
	"cmdb-agent/client/utils"
	"cmdb-agent/daemon/service"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
)

type upgradeCommand struct {
	config     service.Config
	upgradeUrl string
}

func (c *upgradeCommand) run(*kingpin.ParseContext) error {
	fmt.Printf("upgrading the service %s\n", c.config.Name)
	return utils.UpgradeMyself(c.upgradeUrl)
}

func registerUpgrade(cmd *kingpin.CmdClause) {
	c := new(upgradeCommand)
	s := cmd.Command("upgrade", "upgrade the service").
		Action(c.run)

	s.Flag("upgradeUrl", "upgrade url").
		Default("https://image.yeastar.com/tools/cmdb-agent").
		StringVar(&c.upgradeUrl)

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
