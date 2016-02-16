package step
import (
	"github.com/LastCallMedia/vagabond/config"
	"os/exec"
	"fmt"
	"errors"
	"bytes"
)

var ServicesStep = ConfigStep{
	Name: "service containers",
	NeedsRun: func(envt *config.Environment) bool {
		out, err := dockerCommand(envt, "inspect vagabond_proxy").Output()
		if err != nil || bytes.Contains(out, []byte("running")) {
			return true
		}
		out, err = dockerCommand(envt, "inspect vagabond_dnsmasq").Output()
		return err != nil || !bytes.Contains(out, []byte("running"))
	},
	Run: func(envt *config.Environment) error {
		dockerCommand(envt, "stop vagabond_proxy vagabond_dnsmasq").Run()
		dockerCommand(envt, "rm vagabond_proxy vagabond_dnsmasq").Run()
		cmd := dockerCommand(envt, "run --name vagabond_proxy -d -p 80:80 -v /var/run/docker.sock:/tmp/docker.sock:ro jwilder/nginx-proxy")
		if out, err := cmd.CombinedOutput(); err != nil {
			return errors.New(string(out))
		}
		cmd = dockerCommand(envt, fmt.Sprintf("run --name vagabond_dnsmasq -d -p 53:53/udp -p 53:53/tcp --cap-add NET_ADMIN andyshinn/dnsmasq --address=/docker/%s", envt.DockerDaemonIp))
		if out, err := cmd.CombinedOutput(); err != nil {
			return errors.New(string(out))
		}
		return nil
	},
}

func dockerCommand(envt *config.Environment, command string) (cmd *exec.Cmd) {
	if envt.RequiresMachine() {
		cmd = exec.Command("bash", "-c", fmt.Sprintf("docker `docker-machine config %s` %s", envt.MachineName, command))
	} else {
		cmd = exec.Command("bash", "-c", fmt.Sprintf("docker %s", command))
	}
	return
}

