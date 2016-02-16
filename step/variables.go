package step

import (
	"errors"
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/LastCallMedia/vagabond/util"
	"github.com/mitchellh/go-homedir"
	"os/exec"
)

var profileTemplate = `export DOCKER_TZ={{.Tz}}
export VAGABOND_SITES_DIR={{.SitesDir}}
export VAGABOND_DATA_DIR={{.DataDir}}`

var VariablesStep = ConfigStep{
	Name: "environment variables",
	NeedsRun: func(envt *config.Environment) bool {
		profileFilename, err := homedir.Expand("~/.profile")
		if err != nil {
			util.Fatal("Unable to find home directory")
		}
		profile, err := doTemplateAppend(profileTemplate, envt, profileFilename)
		if err != nil {
			return true
		}
		matches, err := checkIfFileMatches(profileFilename, profile)
		if err != nil || !matches {
			return true
		}
		return false
	},
	Run: func(envt *config.Environment) (err error) {
		profileFilename, err := homedir.Expand("~/.profile")
		if err != nil {
			util.Fatal("Unable to find home directory")
		}
		profile, err := doTemplateAppend(profileTemplate, envt, profileFilename)
		if err != nil {
			return
		}
		cmd := exec.Command("tee", profileFilename)
		pipeInputToCmd(cmd, profile)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return errors.New(string(out))
		}
		fmt.Printf(util.FgYellow + "Run the following command once the setup is complete:\n\tsource /etc/profile\n" + util.Reset)
		return
	},
}

