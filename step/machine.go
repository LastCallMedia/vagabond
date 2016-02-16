package step
import (
	"github.com/LastCallMedia/vagabond/config"
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
		}
		newEvt := config.NewEnvironment()
		// Copy over the IPs after the machine boots.
		envt.DockerDaemonIp = newEvt.DockerDaemonIp
		envt.DockerClientIp = newEvt.DockerClientIp
		return
	},
}
