package command

import (
	"cmdb-agent/daemon"
	"context"
	"github.com/alecthomas/kingpin/v2"
	"github.com/drone/signal"
)

type daemonCommand struct {
}

func (c *daemonCommand) run(*kingpin.ParseContext) error {
	ctx, cancel := context.WithCancel(noContext)
	defer cancel()
	// listen for termination signals to gracefully shutdown
	// the runner daemon.
	ctx = signal.WithContextFunc(ctx, func() {
		println("received signal, terminating process")
		cancel()
	})

	return daemon.Run(ctx)
}

func registerDaemon(app *kingpin.Application) {
	c := new(daemonCommand)
	_ = app.Command("daemon", "starts the runner daemon").
		Default().
		Action(c.run)
}
