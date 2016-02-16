package step
import (
	"github.com/LastCallMedia/vagabond/config"
	"fmt"
	"github.com/LastCallMedia/vagabond/util"
)

var MachineStep = ConfigStep{
	Name: "docker machine",
	NeedsRun: func(envt *config.Environment) bool {
		machine := envt.GetMachine()
		return !(machine.IsCreated() && machine.IsBooted())
	},
	Run: func(envt *config.Environment) (err error) {
		machine := envt.GetMachine()
		if ! machine.IsCreated() {
			err = machine.Create().Run()
			if err != nil {
				return
			}
		}
		if ! machine.IsBooted() {
			err = machine.Boot().Run()
			fmt.Printf(util.FgYellow + "Run the following command once the setup is complete:\n\teval $(docker-machine env %s)\n" + util.Reset, envt.MachineName)
		}
		newEvt := config.NewEnvironment()
		// Copy over the IPs after the machine boots.
		envt.DockerDaemonIp = newEvt.DockerDaemonIp
		envt.DockerClientIp = newEvt.DockerClientIp
		return
	},
}
