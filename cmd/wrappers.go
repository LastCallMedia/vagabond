package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
)

var CmdUp = cli.Command{
	Name:  "up",
	Usage: "Start one or more docker containers",
	Action: func(ctx *cli.Context) {
		args := []string{"up", "-d"}
		args = append(args, ctx.Args()...)

		_ = pipeCommand("docker-compose", args...)
	},
}

var CmdDestroy = cli.Command{
	Name:  "destroy",
	Usage: "Remove one or more docker containers",
	Action: func(ctx *cli.Context) {
		args := []string{"rm"}
		args = append(args, ctx.Args()...)
		_ = pipeCommand("docker-compose", args...)
	},
}

var CmdHalt = cli.Command{
	Name:  "halt",
	Usage: "Stop one or more docker containers",
	Action: func(ctx *cli.Context) {
		args := []string{"stop"}
		args = append(args, ctx.Args()...)
		_ = pipeCommand("docker-compose", args...)
	},
}

var CmdSsh = cli.Command{
	Name:  "ssh",
	Usage: "Shell into a running docker container",
	Action: func(ctx *cli.Context) {
		args := []string{"docker", "exec", "-it"}
		for _, i := range ctx.Args() {
			args = append(args, fmt.Sprintf("$(docker-compose ps -q %s)", i))
		}
		args = append(args, "/bin/bash")

		fmt.Printf("running: %s\n", sliceToString(args))
		cmd := exec.Command("bash", "-c", sliceToString(args))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	},
}

func sliceToString(slice []string) string {
	parts := ""
	for _, i := range slice {
		parts += fmt.Sprintf(" %s", i)
	}
	return parts
}

func pipeCommand(name string, arg ...string) error {
	parts := name + sliceToString(arg)
	fmt.Printf("running %s\n", parts)

	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
