package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/LastCallMedia/vagabond/util"
	"github.com/codegangsta/cli"
	"net"
	"os/exec"
	"runtime"
)

// Tests whether setup is completed properly
var CmdDiagnose = cli.Command{
	Name:   "diagnose",
	Usage:  "Run diagnostic tests on the environment",
	Action: runDiagnose,
}

func runDiagnose(ctx *cli.Context) {
	fmt.Println("Running diagnostics...")

	env := config.NewEnvironment()
	if err := env.Check(); err != nil {
		util.Fatal(err)
	}
	if err := checkInstall(env); err != nil {
		util.Fatal(err)
	}
	if err := checkConnection(env); err != nil {
		util.Fatal(err)
	}
	if err := checkContainers(env); err != nil {
		util.Fatal(err)
	}
	if err := checkDns(env); err != nil {
		util.Fatal(err)
	}

	util.Success("OK - No issues found")
}

func checkInstall(env *config.Environment) (err error) {
	err = exec.Command("which", "docker").Run()
	if err != nil {
		return errors.New("docker is not installed. " + helpMissingDocker())
	}
	err = exec.Command("which", "docker-compose").Run()
	if err != nil {
		return errors.New("docker-compose is not installed. " + helpMissingDockerCompose())
	}
	if runtime.GOOS == "darwin" {
		err = exec.Command("which", "docker-machine").Run()
		if err != nil {
			return errors.New("docker-machine is not installed. " + helpMissingDocker())
		}
	}
	return
}

func checkConnection(env *config.Environment) error {
	err := exec.Command("docker", "info").Run()
	if err != nil {
		if runtime.GOOS == "darwin" {
			machine := env.GetMachine()
			if !machine.IsCreated() {
				return errors.New("Docker machine is not created. Run setup to create the machine.")
			}
			if !machine.IsBooted() {
				return errors.New("Docker machine is created but not booted. Run setup to boot the machine.")
			}
		}
		return errors.New("Unable to connect to docker daemon. " + helpConnectingToDaemon(env))
	}
	return err
}

func checkContainers(env *config.Environment) error {
	running := []byte("running")
	out, err := exec.Command("docker", "inspect", "-f", "{{.State.Status}}", "vagabond_dnsmasq").Output()
	if err != nil {
		return errors.New("Problem inspecting DNSMasq container. Run setup to restart it.")
	}
	if !bytes.Contains(out, running) {
		return errors.New("dnsmasq container is not running.  Run setup to restart it.")
	}
	out, err = exec.Command("docker", "inspect", "-f", "{{.State.Status}}", "vagabond_proxy").Output()
	if err != nil {
		return errors.New("Unable to find proxy container. Run setup to start it.")
	}
	if !bytes.Contains(out, running) {
		return errors.New("Proxy container is not running. Run setup to restart it.")
	}
	return err
}

func checkDns(env *config.Environment) error {
	addrs, err := net.LookupIP("somehost.docker")
	if err != nil {
		return errors.New("Unable to resolve somehost.docker. Run setup to fix DNS settings.")
	}

	if !util.IpSliceContains(addrs, env.MachineIp) {
		return errors.New("somehost.docker resolves to the wrong host. Run setup to fix DNS settings.")
	}
	return err
}

func helpConnectingToDaemon(env *config.Environment) string {
	if runtime.GOOS == "darwin" {
		return fmt.Sprintf(`Try running "eval $(docker-machine env %s)"`, env.MachineName)
	}
	return "Make sure the docker service is running and that you are running this command as a user that can access it (usually root)"
}

func helpMissingDocker() string {
	if runtime.GOOS == "darwin" {
		return "Download and install the docker toolbox from https://www.docker.com/products/docker-toolbox"
	}
	return "Install docker using your package manager (https://docs.docker.com/engine/installation/)"
}

func helpMissingDockerCompose() string {
	if runtime.GOOS == "darwin" {
		return "Download and install the docker toolbox from https://www.docker.com/products/docker-toolbox"
	}
	return "Download docker-compose from https://github.com/docker/compose/releases"
}
