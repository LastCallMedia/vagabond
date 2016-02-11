package cmd

import (
	"bufio"
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/LastCallMedia/vagabond/util"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

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

	if env.RequiresMachine() {
		fmt.Println("Ensuring machine is created and booted...")
		err := env.GetMachine().BootOrDie()
		if err != nil {
			util.Fatalf("Unable to boot machine: %s", err)
		}
		// Reset the environment
		env = config.NewEnvironment()
	}

	env.SitesDir = promptQuestion("Enter the sites directory", env.SitesDir)
	env.Tz = promptQuestion("Enter the timezone", env.Tz)
	env.DataDir = promptQuestion("Enter the database storage directory", env.DataDir)

	env.Check()

	coll := []config.ConfigFile{}
	coll = append(coll, config.ProfileConfigFile)
	coll = append(coll, config.ResolverConfigFile)

	if env.RequiresMachine() {
		coll = append(coll, config.NfsExportsConfigFile)
		coll = append(coll, config.BootLocalConfigFile)
	}

	for _, cf := range coll {
		cf.Update(env, ctx.Bool("force"))
	}

	util.Success("All set")
	if env.RequiresMachine() {
		fmt.Printf(`You will also need to run the following commands:
	eval $(docker-machine env %s)
	source /etc/profile
`, env.MachineName)
	}
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
