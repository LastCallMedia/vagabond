package config

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

const (
	vagabond_autoconfig_start = "#VAGABONDAUTOCONFIG"
	vagabond_autoconfig_end   = "#VAGABONDAUTOCONFIGEND"
)

type ConfigFile struct {
	Filename  string
	Generator Generator
	Append    bool
	Io        ConfigFileIo
	Flush     func(env *Environment, filename string)
}

func (cf ConfigFile) Update(env *Environment) {
	existing := cf.Io.Read(env, cf.Filename)
	block := cf.Generator.Generate(env)

	if cf.Append {
		block = appendBlock(existing, block)
	}

	if md5.Sum(existing) != md5.Sum(block) {
		log.Println("Updating " + cf.Filename)
		cf.Io.Write(env, cf.Filename, block)
		// Run any post-update actions.
		cf.Flush(env, cf.Filename)
	}
}

func appendBlock(existing []byte, block []byte) (newblock []byte) {
	re := regexp.MustCompile("(?s)" + vagabond_autoconfig_start + ".*" + vagabond_autoconfig_end)
	newblock = re.ReplaceAll(existing, []byte(""))
	newblock = bytes.TrimRight(newblock, "\n")

	newblock = append(newblock, "\n"+vagabond_autoconfig_start+"\n"...)
	newblock = append(newblock, block...)
	newblock = append(newblock, "\n"+vagabond_autoconfig_end+"\n"...)
	return
}

var NfsExportsConfigFile = ConfigFile{
	Filename: "/etc/exports",
	Io:       ConfigFileLocalIo{},
	Append:   true,
	Generator: Generator{
		Template: `{{.UsersDir}} {{.MachineIp}} -alldirs -mapall=501:20
{{.DataDir}} {{.MachineIp}} -alldirs -maproot=0`,
		TemplateName: "nfsexports",
	},
	Flush: func(env *Environment, filename string) {
		log.Println("Restarting NFS...")
		_, err := exec.Command("sudo", "nfsd", "restart").Output()
		if err != nil {
			log.Fatal("Failed restarting NFS daemon")
		}
	},
}

var ProfileConfigFile = ConfigFile{
	Filename: "/etc/profile",
	Io:       ConfigFileLocalIo{},
	Append:   true,
	Generator: Generator{
		TemplateName: "etcprofile",
		Template: `export DOCKER_TZ={{.Tz}}
export VAGABOND_SITES_DIR={{.SitesDir}}
export VAGABOND_DATA_DIR={{.DataDir}}
export VAGABOND_MACHINE_NAME={{.MachineName}}`,
	},
	Flush: func(env *Environment, filename string) {
		exec.Command("bash", "source", "/etc/profile")
		log.Println("Reload environment variables by running: source /etc/profile")
	},
}

var BootLocalConfigFile = ConfigFile{
	Filename: "/var/lib/boot2docker/bootlocal.sh",
	Io:       ConfigFileMachineIo{},
	Append:   false,
	Generator: Generator{
		TemplateName: "bootlocalsh",
		Template: `#!/bin/sh

sudo umount /Users
	sudo mkdir -p /Users
	sudo /usr/local/etc/init.d/nfs-client start
		sudo mount -t nfs -o noacl,async {{.HostIp}}:{{.UsersDir}} {{.UsersDir}}
		sudo mount -t nfs -o noacl,async {{.HostIp}}:{{.DataDir}} {{.DataDir}}

		docker stop vagabond_proxy vagabond_dnsmasq
		docker rm vagabond_proxy vagabond_dnsmasq
		docker run --name vagabond_proxy   -d -p 80:80 -v /var/run/docker.sock:/tmp/docker.sock:ro jwilder/nginx-proxy
		docker run --name vagabond_dnsmasq -d -p 53:53/udp -p 53:53/tcp --cap-add NET_ADMIN andyshinn/dnsmasq --address=/docker/192.168.99.101
		`,
	},
	Flush: func(env *Environment, filename string) {
		log.Println("Rebooting docker machine...")
		machine := env.GetMachine()
		// Make the script executable. Unless we sync, machine may
		// reboot before the file is written to the disk.
		_, err := machine.Exec(fmt.Sprintf("sudo chmod +x %s", filename)).Output()
		if err != nil {
			log.Fatalf("Failed making %s executable", filename)
		}
		_, err = machine.Reboot().Output()
		if err != nil {
			log.Fatal("Failed rebooting machine")
		}
	},
}

var ResolverConfigFile = ConfigFile{
	Filename: "/etc/resolver/docker",
	Io:       ConfigFileLocalIo{},
	Generator: Generator{
		TemplateName: "etchosts",
		Template:     "nameserver {{.MachineIp}}",
	},
	Flush: func(env *Environment, filename string) {
		// No-op
	},
}
