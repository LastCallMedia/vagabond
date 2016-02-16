package step

import (
	"errors"
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/LastCallMedia/vagabond/util"
	"os/exec"
)

var profileTemplate = `export DOCKER_TZ={{.Tz}}
export VAGABOND_SITES_DIR={{.SitesDir}}
export VAGABOND_DATA_DIR={{.DataDir}}`

var VariablesStep = ConfigStep{
	Name: "environment variables",
	NeedsRun: func(envt *config.Environment) bool {
		profile, err := doTemplateAppend(profileTemplate, envt, "/etc/profile")
		if err != nil {
			return true
		}
		matches, err := checkIfFileMatches("/etc/profile", profile)
		if err != nil || !matches {
			return true
		}
		return false
	},
	Run: func(envt *config.Environment) (err error) {
		profile, err := doTemplateAppend(profileTemplate, envt, "/etc/profile")
		if err != nil {
			return
		}
		fmt.Println("Editing /etc/profile... sudo privileges may be required")
		cmd := exec.Command("sudo", "tee", "/etc/profile")
		pipeInputToCmd(cmd, profile)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return errors.New(string(out))
		}
		fmt.Printf(util.FgYellow + "Run the following command once the setup is complete:\n\tsource /etc/profile\n" + util.Reset)
		return
	},
}
