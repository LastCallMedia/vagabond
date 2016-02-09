package cmd

import (
	"errors"
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/exec"
	"runtime"
)

var CmdDiagnose = cli.Command{
	Name:   "diagnose",
	Action: runDiagnose,
}

func runDiagnose(ctx *cli.Context) {
	fmt.Println("Running diagnostics...")

	env := config.NewEnvironment()
	env.Check()

	err := runDockerInstallationCheck()
	checkOrFatal(err, "Docker is not installed: %s")

	err = runDockerComposeInstallationCheck()
	checkOrFatal(err, "Docker-compose is not installed: %s")

	err = runDockerConnectionCheck()
	if err != nil {
		if runtime.GOOS == "darwin" {
			// Requires machine.
			machine := env.GetMachine()
			if !machine.IsCreated() {
				log.Fatal("Docker machine is not created.")
			}
			if !machine.IsBooted() {
				log.Fatal("Docker machine is created but not booted.")
			}
		}
		log.Fatalf("Docker is unable to connect to the daemon: %s", err)
	}
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
