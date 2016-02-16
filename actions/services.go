package actions
import (
	"github.com/LastCallMedia/vagabond/config"
	"os/exec"
	"fmt"
	"errors"
)

type Services struct {

}

func (act Services)GetName() string {
	return "services"
}

func (act Services)NeedsRun(envt *config.Environment) (bool, error) {
	err := dockerCommand(envt, "inspect vagabond_proxy").Run()
	if err != nil {
		return true, nil
	}
	err = dockerCommand(envt, "inspect vagabond_dnsmasq").Run()
	if err != nil {
		return true, nil
	}
	return false, nil
}

func (act Services)Run(envt *config.Environment) (error) {
	dockerCommand(envt, "stop vagabond_proxy vagabond_dnsmasq").Run()
	dockerCommand(envt, "rm vagabond_proxy vagabond_dnsmasq").Run()
	cmd := dockerCommand(envt, "run --name vagabond_proxy -d -p 80:80 -v /var/run/docker.sock:/tmp/docker.sock:ro jwilder/nginx-proxy")
	if out, err := cmd.Output(); err != nil {
		return errors.New(string(out))
	}
	cmd = dockerCommand(envt, fmt.Sprintf("run --name vagabond_dnsmasq -d -p 53:53/udp -p 53:53/tcp --cap-add NET_ADMIN andyshinn/dnsmasq --address=/docker/%s", envt.MachineIp))
	if out, err := cmd.Output(); err != nil {
		return errors.New(string(out))
	}
	return nil
}

func dockerCommand(envt *config.Environment, command string) (cmd *exec.Cmd) {
	if envt.RequiresMachine() {
		cmd = exec.Command("bash", "-c", fmt.Sprintf("docker `docker-machine config %s` %s", envt.MachineName, command))
	} else {
		cmd = exec.Command("bash", "-c", fmt.Sprintf("docker %s", command))
	}
	return
}

