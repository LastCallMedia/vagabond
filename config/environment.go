package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

const (
	vagabond_tz           string = "America/New_York"
	vagabond_docker_data  string = "/var/lib/dockerdata"
	vagabond_machine_name string = "vagabond"
)

type Environment struct {
	Tz          string
	SitesDir    string
	DataDir     string
	MachineName string
	HostIp      net.IP
	MachineIp   net.IP
	UsersDir    string
}

func NewEnvironment() *Environment {
	var sitesDir, dataDir, tz, machineName string
	var hostIp, machineIp net.IP

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
		dataDir = vagabond_docker_data
	}
	machineName, set = os.LookupEnv("VAGABOND_MACHINE")
	if !set {
		machineName = vagabond_machine_name
	}
	if runtime.GOOS == "darwin" {
		machine := Machine{Name: machineName}
		if machine.IsBooted() {
			hostIp = machine.GetHostIp()
			machineIp = machine.GetIp()
		}
	}
	return &Environment{
		Tz:          tz,
		SitesDir:    sitesDir,
		DataDir:     dataDir,
		MachineName: machineName,
		MachineIp:   machineIp,
		HostIp:      hostIp,
		UsersDir:    "/Users",
	}
}

func (e *Environment) Check() {
	_, err := time.LoadLocation(e.Tz)
	if err != nil {
		log.Fatal("Invalid timezone: ", e.Tz)
	}
	checkDir(e.SitesDir, "Sites directory")
	checkDir(e.DataDir, "Data directory")
}

func checkDir(dir string, name string) {
	src, err := os.Stat(dir)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s does not exist: %s", name, dir))
	}
	if !src.IsDir() {
		log.Fatal(fmt.Sprintf("%s is not a directory: ", name, dir))
	}
}

func (e *Environment) RequiresMachine() bool {
	return runtime.GOOS == "darwin"
}

func (e *Environment) GetMachine() *Machine {
	return &Machine{Name: e.MachineName}
}
