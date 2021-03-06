package cmd

import (
	"bytes"
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/LastCallMedia/vagabond/util"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
)

// Proxies a docker-compose up command.
var CmdUp = cli.Command{
	Name:  "up",
	Usage: "Start one or more docker containers",
	Action: func(ctx *cli.Context) {
		runAndNotifyCommand("docker-compose", append([]string{"up", "-d"}, ctx.Args()...)...)
	},
}

// Proxies a docker-compose rm command.
var CmdDestroy = cli.Command{
	Name:            "destroy",
	Aliases:         []string{"rm"},
	Usage:           "Remove one or more docker containers",
	SkipFlagParsing: true,
	Action: func(ctx *cli.Context) {
		runAndNotifyCommand("docker-compose", append([]string{"rm"}, ctx.Args()...)...)
	},
}

// Proxies a docker-compose stop command.
var CmdHalt = cli.Command{
	Name:    "halt",
	Aliases: []string{"stop"},
	Usage:   "Stop one or more docker containers",
	Action: func(ctx *cli.Context) {
		runAndNotifyCommand("docker-compose", append([]string{"stop"}, ctx.Args()...)...)
	},
}

// Proxies a docker-compose ps command.
var CmdStatus = cli.Command{
	Name:    "status",
	Aliases: []string{"ps"},
	Usage:   "View the status of running containers",
	Action: func(ctx *cli.Context) {
		runAndNotifyCommand("docker-compose", append([]string{"ps"}, ctx.Args()...)...)
	},
}

// Proxies a docker-machine ip command.
var CmdIp = cli.Command{
	Name:  "ip",
	Usage: "View the IP Address of the docker machine",
	Action: func(ctx *cli.Context) {
		env := config.NewEnvironment()
		runAndNotifyCommand("docker-machine", "ip", env.MachineName)
	},
}

// Emulates vagrant ssh with docker exec
var CmdSsh = cli.Command{
	Name:    "ssh",
	Aliases: []string{"exec"},
	Usage:   "Shell into a running docker container",
	Action: func(ctx *cli.Context) {
		numArgs := len(ctx.Args())
		if numArgs > 1 {
			notifyError("You may only specify a single container")
			os.Exit(1)
		}
		if numArgs < 1 {
			notifyError("You must specify a container name")
			os.Exit(1)
		}
		name := ctx.Args()[0]
		contid, err := exec.Command("docker-compose", "ps", "-q", name).Output()
		contid = bytes.TrimSpace(contid)
		if err != nil || bytes.Equal(contid, []byte("")) {
			fmt.Printf("Could not find container %s.  Are you sure it's running?\n", name)
			os.Exit(1)
		}

		notifyCommand("docker", "exec", "-it", fmt.Sprintf("$(docker-compose ps -q %s)", name), "/bin/bash")
		_ = pipeCommand("docker", "exec", "-it", string(contid), "/bin/bash")
	},
}

func sliceToString(slice []string) string {
	parts := ""
	for _, i := range slice {
		parts += fmt.Sprintf(" %s", i)
	}
	return parts
}

func runAndNotifyCommand(name string, arg ...string) {
	notifyCommand(name, arg...)
	pipeCommand(name, arg...)
}

func pipeCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func notifyCommand(name string, arg ...string) {
	parts := name + sliceToString(arg)
	fmt.Printf("%sRunning: %s%s%s\n", util.Dim, util.Reset+util.Bright+util.FgGreen, parts, util.Reset)
}

func notifyError(text string) {
	fmt.Printf("%s%s%s\n", util.FgRed, text, util.Reset)
}
