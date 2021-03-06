package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

const (
	vagabond_tz           string = "America/New_York"
	vagabond_machine_name string = "vagabond"
)

// Representation of the vagabond environment settings.
type Environment struct {
	Tz             string
	SitesDir       string
	DataDir        string
	MachineName    string
	DockerClientIp net.IP
	DockerDaemonIp net.IP
	UsersDir       string
}

// Create and prepopulate a new environment based on settings
func NewEnvironment() *Environment {
	var sitesDir, dataDir, tz, machineName string
	clientIp := net.ParseIP("127.0.0.1")
	daemonIp := net.ParseIP("127.0.0.1")

	tz, set := os.LookupEnv("DOCKER_TZ")
	if !set {
		tz = vagabond_tz
	}
	sitesDir, set = os.LookupEnv("VAGABOND_SITES_DIR")
	if !set {
		if runtime.GOOS == "darwin" {
			sitesDir = os.ExpandEnv("$HOME/Sites")
		} else {
			sitesDir = "/var/www"
		}
	}

	dataDir, set = os.LookupEnv("VAGABOND_DATA_DIR")
	if !set {
		if runtime.GOOS == "darwin" {
			dataDir = os.ExpandEnv("$HOME/Library/Vagabond")
		} else {
			dataDir = "/var/lib/dockerdata"
		}
	}
	if runtime.GOOS == "darwin" {
		machineName = vagabond_machine_name
		machine := Machine{Name: machineName}
		if machine.IsBooted() {
			clientIp = machine.GetHostIp()
			daemonIp = machine.GetIp()
		}
	}

	return &Environment{
		Tz:             tz,
		SitesDir:       sitesDir,
		DataDir:        dataDir,
		MachineName:    machineName,
		DockerDaemonIp: daemonIp,
		DockerClientIp: clientIp,
		UsersDir:       "/Users",
	}
}

// Verify that environment variables are set properly
func (e *Environment) Check() error {
	if _, err := time.LoadLocation(e.Tz); err != nil {
		return errors.New(fmt.Sprintf("Invalid timezone: %s", e.Tz))
	}
	if err := checkDir(e.SitesDir, "Sites directory"); err != nil {
		return err
	}
	if err := checkDir(e.DataDir, "Data directory"); err != nil {
		return err
	}
	return nil
}

func checkDir(dir string, name string) error {
	src, err := os.Stat(dir)
	if err != nil {
		return errors.New(fmt.Sprintf("%s does not exist: %s", name, dir))
	}
	if !src.IsDir() {
		return errors.New(fmt.Sprintf("%s is not a directory: ", name, dir))
	}
	return nil
}

// Assert whether the environment requires docker machine to run
func (e *Environment) RequiresMachine() bool {
	return runtime.GOOS == "darwin"
}

// Get the docker machine instance for the environment.
func (e *Environment) GetMachine() *Machine {
	return &Machine{Name: e.MachineName}
}
