package cmd

import (
	"bufio"
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"strings"
)

var CmdConfigure = cli.Command{
	Name:   "configure",
	Action: runConfigure,
}

func runConfigure(ctx *cli.Context) {
	env := config.NewEnvironment()

	if env.RequiresMachine() {
		machine := env.GetMachine()
		if !machine.IsCreated() {
			log.Println("Creating the machine...")
			_, err := machine.Create().Output()
			if err != nil {
				log.Fatal("Error creating machine: ", err)
			}
			// Reset the environment
			env = config.NewEnvironment()
		}
		if !machine.IsBooted() {
			log.Println("Booting the machine...")
			_, err := machine.Boot().Output()
			if err != nil {
				log.Fatal("Error booting machine: ", err)
			}
			// Reset the environment
			env = config.NewEnvironment()
		}
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
		cf.Update(env)
	}

	if env.RequiresMachine() {
		fmt.Printf(`
All set. You will also need to run the following commands:
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
		log.Fatal("Error reading input")
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return def
	}
	return input
}
