package main

import (
	"github.com/LastCallMedia/vagabond/cmd"
	"github.com/codegangsta/cli"
	"os"
)

// Version tag.  This gets replaced when compiling a tag build.
var Version = "development"

func main() {
	app := cli.NewApp()
	app.Name = "Vagabond"
	app.Usage = "Development environment helper"
	app.Version = Version
	app.Commands = []cli.Command{
		cmd.CmdSetup,
		cmd.CmdDiagnose,
		cmd.CmdSelfUpdate,

		cmd.CmdUp,
		cmd.CmdDestroy,
		cmd.CmdHalt,
		cmd.CmdSsh,
	}
	app.Run(os.Args)
}
