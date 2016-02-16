package actions
import (
	"github.com/LastCallMedia/vagabond/config"
	"os/exec"
	"errors"
)

var VariablesTemplate = `export DOCKER_TZ={{.Tz}}
export VAGABOND_SITES_DIR={{.SitesDir}}
export VAGABOND_DATA_DIR={{.DataDir}}`

type VariablesAction struct {
	
}

func (act VariablesAction)GetName() string {
	return "environment setup"
}

func (act VariablesAction)NeedsRun(envt *config.Environment) (bool, error) {
	profile, err := doTemplateAppend(VariablesTemplate, envt, "/etc/profile")
	if err != nil {
		return true, nil
	}
	matches, err := checkIfFileMatches("/etc/profile", profile)
	if err != nil || !matches {
		return true, nil
	}
	return false, nil
}

func (act VariablesAction)Run(envt *config.Environment) (err error) {
	profile, err := doTemplateAppend(VariablesTemplate, envt, "/etc/profile")
	if err != nil {
		return
	}

	cmd := exec.Command("sudo", "tee", "/etc/profile")
	pipeInputToCmd(cmd, profile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return
}