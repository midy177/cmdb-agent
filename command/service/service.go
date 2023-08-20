package service

import (
	"github.com/alecthomas/kingpin/v2"
)

const (
	// DefaultName is the default service name.
	DefaultName = "cmdb-agent"

	// DefaultDesc is the default service description.
	DefaultDesc = "cmdb agent"
)

// Register registers the command.
func Register(app *kingpin.Application) {
	cmd := app.Command("service", "manages the runner service")
	registerInstall(cmd)
	registerStart(cmd)
	registerStop(cmd)
	registerUninstall(cmd)
	registerRun(cmd)
	registerRestart(cmd)
}
