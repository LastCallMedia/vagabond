package cmd

import (
	"bufio"
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/LastCallMedia/vagabond/util"
	"github.com/codegangsta/cli"
	"os"
	"strings"
	"github.com/LastCallMedia/vagabond/step"
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
	env := config.NewEnvironment()
	force := ctx.Bool("force")

	env.SitesDir = promptQuestion("Enter the sites directory", env.SitesDir)
	env.Tz = promptQuestion("Enter the timezone", env.Tz)
	env.DataDir = promptQuestion("Enter the database storage directory", env.DataDir)

	env.Check()

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

func promptQuestion(question string, def string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s[%s]: ", question, def)
	input, err := reader.ReadString('\n')

	if err != nil {
		util.Fatal("Error reading input")
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return def
	}
	return input
}
