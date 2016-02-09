package cmd

import (
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/codegangsta/cli"
	"log"
	"os/exec"
	"runtime"
	"errors"
)

var DockerInstallHelp = `Download and install the docker toolbox from https://www.docker.com/products/docker-toolbox`
var DockerComposeInstallHelp = DockerInstallHelp

var CmdDiagnose = cli.Command{
	Name:   "diagnose",
	Action: runDiagnose,
}

func runDiagnose(ctx *cli.Context) {
	fmt.Println("Running diagnostics...")

	env := config.NewEnvironment()
	env.Check()

	err := checkInstall(env)
	if err != nil {
		log.Fatal(err)
	}

	err = checkConnection(env)
	if err != nil {
		log.Fatal(err)
	}
}

func checkInstall(env *config.Environment) (err error) {
	err = exec.Command("which", "docker").Run()
	if err != nil {
		return errors.New("docker is not installed. " + DockerInstallHelp)
	}
	err = exec.Command("which", "docker-compose").Run()
	if err != nil {
		return errors.New("docker-compose is not installed. " + DockerComposeInstallHelp)
	}
	return
}

func checkConnection(env *config.Environment) (err error) {
	err = exec.Command("docker", "info").Run()
	if err != nil {
		if runtime.GOOS == "darwin" {
			machine := env.GetMachine()
			if !machine.IsCreated() {
				return errors.New("Docker machine is not created. Run configure to create the machine.")
			}
			if !machine.IsBooted() {
				return errors.New("Docker machine is created but not booted. Run configure to boot the machine.")
			}
		}
		return errors.New("Docker daemon is not running.")
	}
	return
}
