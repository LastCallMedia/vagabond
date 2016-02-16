package actions

import(
	"github.com/LastCallMedia/vagabond/config"
//	"net"
	"os/exec"
	"errors"
	"fmt"
)

type DnsAction struct {
}

func (act DnsAction) GetName() string {
	return "DNS client"
}

type DnsActionOsx struct {
	DnsAction
}

func (act DnsActionOsx)NeedsRun(envt *config.Environment) (bool, error) {
//	_, err := net.LookupIP("ping.docker")
//	if err != nil {
//		return true, nil
//	}
	return true, nil
}

func (act DnsActionOsx)Run(envt *config.Environment) (error) {
	err := exec.Command("sudo", "mkdir", "-p", "/etc/resolver").Run()
	if err != nil {
		return errors.New("Unable to create /etc/resolver")
	}
	cmd := exec.Command("sudo", "tee", "/etc/resolver/docker")
	contents := []byte(fmt.Sprintf("nameserver %s", envt.MachineIp))
	pipeInputToCmd(cmd, contents)
	err = cmd.Run()
	if err != nil {
		return errors.New("Unable to write /etc/resolver/docker")
	}

	return nil
}

type DnsActionLinux struct {
	DnsAction
}

func (act DnsActionLinux)NeedsRun(envt *config.Environment) (bool, error) {
	return true, nil
}
func (act DnsActionLinux)Run(envt *config.Environment) (error) {
	return nil
}

