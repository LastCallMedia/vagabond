package cmd

import (
	"errors"
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/exec"
)

var CmdDiagnose = cli.Command{
	Name:   "diagnose",
	Action: runDiagnose,
}

func runDiagnose(ctx *cli.Context) {
	fmt.Println("Running diagnostics...")

	err := runDockerInstallationCheck()
	checkOrFatal(err, "Docker is not installed: %s")

	err = runDockerComposeInstallationCheck()
	checkOrFatal(err, "Docker-compose is not installed: %s")

	err = runDockerConnectionCheck()
	checkOrFatal(err, "Docker is unable to connect to the daemon: %s")

	env := config.NewEnvironment()
	env.Check()
}

func checkOrFatal(err error, msg string) {
	if err != nil {
		log.Fatal(fmt.Sprintf(msg, err))
	}
}

func runDockerInstallationCheck() (err error) {
	cmd := exec.Command("which", "docker")
	return cmd.Run()
}

func runDockerComposeInstallationCheck() (err error) {
	cmd := exec.Command("which", "docker-compose")
	return cmd.Run()
}

func runDockerConnectionCheck() (err error) {
	cmd := exec.Command("docker", "info")
	return cmd.Run()
}

func runEnvCheck() (err error) {
	envvar := os.Getenv("DOCKER_SITES_DIR")
	if "" == envvar {
		err = errors.New("$DOCKER_SITES_DIR is not defined")
		return err
	}

	envvar = os.Getenv("DOCKER_TZ")
	if "" == envvar {
		err = errors.New("$DOCKER_TZ is not defined")
		return err
	}

	return err
}
