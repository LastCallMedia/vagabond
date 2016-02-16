package step
import(
	"github.com/LastCallMedia/vagabond/config"
	"fmt"
	"errors"
)

var BootLocalTemplate = `#!/bin/bash
	sudo mkdir -p /Users
	sudo mkdir -p {{.DataDir}}
	sudo /usr/local/etc/init.d/nfs-client start
		sudo mount -t nfs -o noacl,async {{.DockerClientIp}}:{{.UsersDir}} {{.UsersDir}}
		sudo mount -t nfs -o noacl,async {{.DockerClientIp}}:{{.DataDir}} {{.DataDir}}
`

var NfsClientStep = ConfigStep{
	Name: "nfs client",
	NeedsRun: func(envt *config.Environment) bool {
		machine := envt.GetMachine()
		mountErr := machine.Exec(fmt.Sprintf("mount -t nfs|grep %s", envt.UsersDir)).Run()
		if mountErr != nil {
			return true
		}
		mountErr = machine.Exec(fmt.Sprintf("mount -t nfs|grep %s", envt.DataDir)).Run()

		return mountErr != nil
	},
	Run: func(envt *config.Environment) (err error) {
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
	},
}
