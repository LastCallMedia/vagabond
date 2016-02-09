package config

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

type Machine struct {
	Name string
}

func NewMachine(name string) *Machine {
	return &Machine{Name: name}
}

func (m *Machine) IsCreated() bool {
	_, err := exec.Command("docker-machine", "status", m.Name).Output()
	return err == nil
}

func (m *Machine) IsBooted() bool {
	out, err := exec.Command("docker-machine", "status", m.Name).Output()
	return err == nil && bytes.Contains(out, []byte("Running"))
}

func (m *Machine) BootOrDie() (err error) {
	if !m.IsCreated() {
		_, err = m.Create().Output()
		if err != nil {
			return
		}
	}
	if !m.IsBooted() {
		_, err = m.Boot().Output()
		if err != nil {
			return
		}
	}
	return
}

func (m *Machine) Exec(cmd string) *exec.Cmd {
	return exec.Command("docker-machine", "ssh", m.Name, cmd)
}

func (m *Machine) Scp(localFile string, remoteFile string) *exec.Cmd {
	return exec.Command("docker-machine", "scp", localFile, fmt.Sprintf("%s:%s", m.Name, remoteFile))
}

func (m *Machine) Create() *exec.Cmd {
	return exec.Command("docker-machine", "create", m.Name, "-d", "virtualbox")
}

func (m *Machine) Reboot() *exec.Cmd {
	return exec.Command("docker-machine", "restart", m.Name)
}

func (m *Machine) Boot() *exec.Cmd {
	return exec.Command("docker-machine", "start", m.Name)
}

func (m *Machine) GetIp() net.IP {
	out, err := exec.Command("docker-machine", "ip", m.Name).Output()
	if err != nil {
		log.Fatal(err)
	}

	out = bytes.TrimSpace(out)
	return net.ParseIP(string(out))
}

func (m *Machine) GetHostIp() net.IP {
	out, err := exec.Command("docker-machine", "inspect", m.Name, "-f", "{{.Driver.HostOnlyCIDR}}").Output()
	if err != nil {
		log.Fatal(err)
	}
	hostCidr := strings.TrimSpace(string(out))
	ip, _, _ := net.ParseCIDR(hostCidr)

	return ip
}
