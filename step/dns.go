package step

import(
	"github.com/LastCallMedia/vagabond/config"
	"runtime"
	"net"
	"os/exec"
	"errors"
	"fmt"
	"github.com/LastCallMedia/vagabond/util"
)


var dnsResolverStep = ConfigStep{
	Name:"DNS",
	NeedsRun: dnsNeedsConfigure,
	Run: func(envt *config.Environment) (err error) {
		err = exec.Command("sudo", "mkdir", "-p", "/etc/resolver").Run()
		if err != nil {
			return errors.New("Unable to create /etc/resolver")
		}
		cmd := exec.Command("sudo", "tee", "/etc/resolver/docker")
		contents := []byte(fmt.Sprintf("nameserver %s", envt.DockerDaemonIp))
		pipeInputToCmd(cmd, contents)
		err = cmd.Run()
		if err != nil {
			return errors.New("Unable to write /etc/resolver/docker")
		}

		return nil
	},
}

var dnsOtherStep = ConfigStep{
	Name: "DNS",
	NeedsRun: dnsNeedsConfigure,
	Run: func(envt *config.Environment) (err error) {
		fmt.Printf(util.FgYellow + "Unable to do automatic configuration of *.docker domains.  Please point your DNS for these domains to %s\n" + util.Reset, envt.DockerClientIp)
		return
	},
}


func dnsNeedsConfigure(envt *config.Environment) bool {
	addrs, err := net.LookupIP("ping.docker")
	return err != nil || !util.IpSliceContains(addrs, envt.DockerDaemonIp)
}

func NewDnsAction() ConfigStep {
	switch runtime.GOOS {
	case "darwin":
		return dnsResolverStep
	default:
		return dnsOtherStep
	}
}


