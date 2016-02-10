package cmd

import (
	"bytes"
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
		name := ctx.Args()[0]
		contid, err := exec.Command("docker-compose", "ps", "-q", name).Output()
		contid = bytes.TrimSpace(contid)
		if err != nil || bytes.Equal(contid, []byte("")) {
			fmt.Printf("Could not find container %s.  Are you sure it's running?\n", name)
			os.Exit(1)
		}

		fmt.Printf("running: docker exec -it $(docker-compose ps -q %s) /bin/bash \n", name)
		cmd := exec.Command("docker", "exec", "-it", string(contid), "/bin/bash")
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
