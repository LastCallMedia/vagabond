package actions
import (
	"github.com/LastCallMedia/vagabond/config"
)

type MachineBootAction struct {

}
func (act MachineBootAction)GetName() string {
	return "boot docker-machine"
}
func (act MachineBootAction)NeedsRun(envt *config.Environment) (bool, error) {
	machine := envt.GetMachine()
	return !(machine.IsCreated() && machine.IsBooted()), nil
}
func (act MachineBootAction)Run(envt *config.Environment) (err error) {
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
	envt.MachineIp = newEvt.MachineIp
	envt.HostIp = newEvt.HostIp
	return
}
