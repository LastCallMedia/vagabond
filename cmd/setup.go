package cmd

import (
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/LastCallMedia/vagabond/util"
	"github.com/codegangsta/cli"
	"os"
	"github.com/LastCallMedia/vagabond/step"
	"github.com/Songmu/prompter"
	"errors"
)

// Sets up the vagabond environment
var CmdSetup = cli.Command{
	Name:        "setup",
	Usage:       "Prepare the environment for use",
	Description: "Prepare the environment for the first use",
	Action:      runSetup,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force",
			Usage: "Force the setup actions",
		},
	},
}

func runSetup(ctx *cli.Context) {
	var err error

	env := config.NewEnvironment()
	force := ctx.Bool("force")

	env.Tz = prompter.Prompt("Timezone", env.Tz)
	env.SitesDir, err = promptForDir("Sites directory", env.SitesDir)
	if err != nil {
		util.Fatal(err)
	}
	env.DataDir, err = promptForDir("Database storage directory", env.DataDir)
	if err != nil {
		util.Fatal(err)
	}

	err = env.Check()
	if err != nil {
		util.Fatal(err)
	}

	acts := []step.ConfigStep{
		step.VariablesStep,
	}

	if env.RequiresMachine() {
		acts = append(acts, step.MachineStep)
		acts = append(acts, step.NfsServerStep)
		acts = append(acts, step.NfsClientStep)
	}
	acts = append(acts, step.ServicesStep)
	acts = append(acts, step.NewDnsAction())


	for _, act := range acts {
		needs := act.NeedsRun(env)
		if needs || force {
			util.Successf("Running %s", act.GetName())
			err := act.Run(env)
			if err != nil {
				util.Fatal(err)
			}
		} else {
			util.Successf("Skipping %s", act.GetName())
		}
	}

	util.Success("Setup complete")
}

func promptForDir(name string, def string) (dir string, err error) {
	dir = prompter.Prompt(name, def)
	exists, err  := util.DirExists(dir)
	if ! exists {
		create := prompter.YN(fmt.Sprintf("%s does not exist.  Create?", dir), false)
		if create {
			err = os.MkdirAll(dir, 0755)
			return
		}
		err = errors.New("Did not create directory")
	}
	return
}

