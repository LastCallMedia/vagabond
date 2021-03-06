package step

import (
	"errors"
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"os/exec"
)

var ExportsTemplate = `{{.UsersDir}} {{.DockerDaemonIp}} -alldirs -mapall=501:20
{{.DataDir}} {{.DockerDaemonIp}} -alldirs -maproot=0`

var NfsServerStep = ConfigStep{
	Name: "nfs server",
	NeedsRun: func(envt *config.Environment) bool {
		exports, err := doTemplateAppend(ExportsTemplate, envt, "/etc/exports")
		if err != nil {
			return true
		}
		matches, err := checkIfFileMatches("/etc/exports", exports)

		return err != nil || !matches
	},
	Run: func(envt *config.Environment) (err error) {
		fmt.Println("Editing /etc/exports... sudo privileges may be required")
		out, err := exec.Command("sudo", "touch", "/etc/exports").CombinedOutput()
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to create/update exports file: %s", string(out)))
		}

		exports, err := doTemplateAppend(ExportsTemplate, envt, "/etc/exports")
		if err != nil {
			return err
		}

		cmd := exec.Command("sudo", "tee", "/etc/exports")
		pipeInputToCmd(cmd, exports)
		out, err = cmd.CombinedOutput()
		if err != nil {
			return errors.New(string(out))
		}

		cmd = exec.Command("sudo", "nfsd", "restart")
		out, err = cmd.CombinedOutput()
		if err != nil {
			return errors.New(string(out))
		}
		return
	},
}
