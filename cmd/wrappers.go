package cmd

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
)

const(
	// Courtesy of https://groups.google.com/forum/#!topic/golang-nuts/99MKtEkvQ2c
	Reset = "\x1b[0m"
	Bright = "\x1b[1m"
	Dim = "\x1b[2m"
	Underscore = "\x1b[4m"
	Blink = "\x1b[5m"
	Reverse = "\x1b[7m"
	Hidden = "\x1b[8m"

	FgBlack = "\x1b[30m"
	FgRed = "\x1b[31m"
	FgGreen = "\x1b[32m"
	FgYellow = "\x1b[33m"
	FgBlue = "\x1b[34m"
	FgMagenta = "\x1b[35m"
	FgCyan = "\x1b[36m"
	FgWhite = "\x1b[37m"
)

var CmdUp = cli.Command{
	Name:  "up",
	Usage: "Start one or more docker containers",
	Action: func(ctx *cli.Context) {
		args := append([]string{"-d"}, ctx.Args()...)
		runSimpleComposeCommand("up", args...)
	},
}

var CmdDestroy = cli.Command{
	Name:  "destroy",
	Aliases: []string{"rm"},
	Usage: "Remove one or more docker containers",
	SkipFlagParsing: true,
	Action: func(ctx *cli.Context) {
		runSimpleComposeCommand("rm", ctx.Args()...)
	},
}

var CmdHalt = cli.Command{
	Name:  "halt",
	Aliases: []string{"stop"},
	Usage: "Stop one or more docker containers",
	Action: func(ctx *cli.Context) {
		runSimpleComposeCommand("stop", ctx.Args()...)
	},
}

var CmdStatus = cli.Command{
	Name: "status",
	Aliases: []string{"ps"},
	Usage: "View the status of running containers",
	Action: func(ctx *cli.Context) {
		runSimpleComposeCommand("ps", ctx.Args()...)
	},
}

var CmdSsh = cli.Command{
	Name:  "ssh",
	Aliases: []string{"exec"},
	Usage: "Shell into a running docker container",
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

func runSimpleComposeCommand(name string, arg ...string) {
	arg = append([]string{name}, arg...)
	notifyCommand("docker-compose", arg...)
	pipeCommand("docker-compose", arg...)
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
	fmt.Printf("%sRunning: %s%s%s\n", Dim, Reset +Bright + FgGreen, parts, Reset)
}

func notifyError(text string) {
	fmt.Printf("%s%s%s\n", FgRed, text, Reset)
}