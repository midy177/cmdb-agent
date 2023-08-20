package command

import (
	"cmdb-agent/command/service"
	"context"
	"github.com/alecthomas/kingpin/v2"
	"os"
)

// program version
var version = "0.0.1"

// empty context
var noContext = context.Background()

// Command parses the command line arguments and then executes a
// subcommand program.
func Command() {
	app := kingpin.New(service.DefaultName, service.DefaultDesc)
	registerDaemon(app)
	service.Register(app)
	kingpin.Version(version)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
