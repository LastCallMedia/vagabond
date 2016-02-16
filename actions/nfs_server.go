package actions

import(
	"github.com/LastCallMedia/vagabond/config"
	"os/exec"
	"errors"
	"fmt"
)

var ExportsTemplate = `{{.UsersDir}} {{.MachineIp}} -alldirs -mapall=501:20
{{.DataDir}} {{.MachineIp}} -alldirs -maproot=0`

type NfsServerAction struct {

}

func (act NfsServerAction)GetName() string {
	return "nfs server"
}

func (act NfsServerAction)NeedsRun(envt *config.Environment) (bool, error) {
	exports, err := doTemplateAppend(ExportsTemplate, envt, "/etc/exports")
	if err != nil {
		return true, nil
	}
	matches, err := checkIfFileMatches("/etc/exports", exports)
	if err != nil || !matches {
		return true, nil
	}
	return false, nil
}

func (act NfsServerAction)Run(envt *config.Environment) (err error) {
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
}

