package main

import (
	"github.com/LastCallMedia/vagabond/cmd"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Vagabond"
	app.Usage = "Vagabond instructions..."
	// @todo: Add versioning, like
	// https://ariejan.net/2015/10/12/building-golang-cli-tools-update/
	app.Version = "1.0"
	app.Commands = []cli.Command{
		cmd.CmdSetup,
		cmd.CmdDiagnose,
	}
	app.Run(os.Args)
}
