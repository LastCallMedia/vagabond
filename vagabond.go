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
	app.Version = "1.0"
	app.Commands = []cli.Command{
		cmd.CmdDiagnose,
		cmd.CmdConfigure,
	}
	app.Run(os.Args)
}
