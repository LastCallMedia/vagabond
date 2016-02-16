package actions
import(
	"github.com/LastCallMedia/vagabond/config"
	"fmt"
	"errors"
)

var BootLocalTemplate = `sudo umount /Users
	sudo mkdir -p /Users
	sudo mkdir -p {{.DataDir}}
	sudo /usr/local/etc/init.d/nfs-client start
		sudo mount -t nfs -o noacl,async {{.HostIp}}:{{.UsersDir}} {{.UsersDir}}
		sudo mount -t nfs -o noacl,async {{.HostIp}}:{{.DataDir}} {{.DataDir}}
`


type NfsClientAction struct {

}

func (act NfsClientAction)GetName() string {
	return "nfs client"
}

func (act NfsClientAction)NeedsRun(envt *config.Environment) (bool, error) {
	machine := envt.GetMachine()
	mountErr := machine.Exec(fmt.Sprintf("mount -t nfs|grep %s", envt.UsersDir)).Run()
	if mountErr != nil {
		return true, nil
	}
	mountErr = machine.Exec(fmt.Sprintf("mount -t nfs|grep %s", envt.DataDir)).Run()
	if mountErr != nil {
		return true, nil
	}
	return false, nil
}

func (act NfsClientAction)Run(envt *config.Environment) (err error) {
	bootLocal, err := doTemplate(BootLocalTemplate, envt)
	if err != nil {
		return
	}
	machine := envt.GetMachine()
	cmd := machine.Exec("sudo tee /var/lib/boot2docker/bootlocal.sh")
	pipeInputToCmd(cmd, bootLocal)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	out, err = machine.Exec("sync; sudo chmod +x /var/lib/boot2docker/bootlocal.sh; sync").CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	err = machine.Reboot().Run()
	return
}