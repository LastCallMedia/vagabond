package cmd

import (
	"bufio"
	"fmt"
	"github.com/LastCallMedia/vagabond/config"
	"github.com/LastCallMedia/vagabond/util"
	"github.com/codegangsta/cli"
	"os"
	"strings"
	"github.com/LastCallMedia/vagabond/cfg"
	"text/template"
)

// Sets up the vagabond environment
var CmdSetup = cli.Command{
	Name:        "setup",
	Usage:       "Prepare the environment for use",
	Description: "Prepare the environment for the first use",
	Action:      runSetup,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force",
			Usage: "Force the setup actions",
		},
	},
}

func runSetup(ctx *cli.Context) {
	env := config.NewEnvironment()
	force := ctx.Bool("force")

	if env.RequiresMachine() {
		fmt.Println("Ensuring machine is created and booted...")
		err := env.GetMachine().BootOrDie()
		if err != nil {
			util.Fatalf("Unable to boot machine: %s", err)
		}
		// Reset the environment
		env = config.NewEnvironment()
	}

	env.SitesDir = promptQuestion("Enter the sites directory", env.SitesDir)
	env.Tz = promptQuestion("Enter the timezone", env.Tz)
	env.DataDir = promptQuestion("Enter the database storage directory", env.DataDir)

	env.Check()

	actions := []cfg.ConfigAction{}

	if env.RequiresMachine() {
		actions = append(actions, nfsExportsActions(env)...)
		actions = append(actions, bootLocalActions(env)...)
		// On osx, use resolver files.
		actions = append(actions, resolverActions(env)...)
	} else {
		// On linux, use dhclient
//		coll = append(coll, config.DhclientConfigFile)
	}
	actions = append(actions, etcProfileAction(env))
	actions = append(actions, containerServiceActions(env)...)


	for _, action := range actions {
		needs, err := action.NeedsRun()
		if err != nil {
			util.Fatalf("Got error checking if %s needs to be run: %s", action.GetName(), err)
		}
		if needs || force {
			util.Successf("Running %s", action.GetName())
			err = action.Run()
			if err != nil {
				util.Fatal(err)
			}
		}
	}

	util.Success("All set")
	if env.RequiresMachine() {
		fmt.Printf(`You will also need to run the following commands:
	eval $(docker-machine env %s)
	source /etc/profile
`, env.MachineName)
	}
}

func promptQuestion(question string, def string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s[%s]: ", question, def)
	input, err := reader.ReadString('\n')

	if err != nil {
		util.Fatal("Error reading input")
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return def
	}
	return input
}

func etcProfileAction(env *config.Environment) (cfg.ConfigAction) {
	return cfg.ConfigFileAction{
		Filename: "/etc/profile",
		Append: true,
		TemplateVars: env,
		Template: template.Must(template.New("etcprofile").Parse(`export DOCKER_TZ={{.Tz}}
export VAGABOND_SITES_DIR={{.SitesDir}}
export VAGABOND_DATA_DIR={{.DataDir}}`)),

	}
}

func nfsExportsActions(env *config.Environment) ([]cfg.ConfigAction) {
	return []cfg.ConfigAction{
		cfg.ConfigFileAction{
			Filename: "/etc/exports",
			Append:   true,
			TemplateVars: env,
			Template: template.Must(template.New("etcexports").Parse(`{{.UsersDir}} {{.MachineIp}} -alldirs -mapall=501:20
{{.DataDir}} {{.MachineIp}} -alldirs -maproot=0`)),
		},
		cfg.CommandConfigAction{
			Command: "sudo nfsd restart",
		},
	}
}

func bootLocalActions(env *config.Environment) ([]cfg.ConfigAction) {
	return []cfg.ConfigAction{
		cfg.ConfigFileAction{
			Filename: "/var/lib/boot2docker/bootlocal.sh",
			Target: env.MachineName,
			TemplateVars: env,
			Template: template.Must(template.New("bootlocalsh").Parse(`#!/bin/sh

sudo umount /Users
	sudo mkdir -p /Users
	sudo /usr/local/etc/init.d/nfs-client start
		sudo mount -t nfs -o noacl,async {{.HostIp}}:{{.UsersDir}} {{.UsersDir}}
		sudo mount -t nfs -o noacl,async {{.HostIp}}:{{.DataDir}} {{.DataDir}}`)),
		},
		cfg.CommandConfigAction{
			Command: "sudo chmod +x /var/lib/boot2docker/bootlocal.sh",
			Target: env.MachineName,
		},
		cfg.CommandConfigAction{
			Command: fmt.Sprintf("docker-machine restart %s", env.MachineName),
		},
	}
}

func containerServiceActions(env *config.Environment) ([]cfg.ConfigAction) {
	return []cfg.ConfigAction{
		cfg.CommandConfigAction{
			Target:env.MachineName, // will be "" on linux
			Command: "docker stop vagabond_proxy vagabond_dnsmasq",
			IgnoreReturn: true,
		},
		cfg.CommandConfigAction{
			Target:env.MachineName, // will be "" on linux
			Command: "docker rm vagabond_proxy vagabond_dnsmasq",
			IgnoreReturn: true,
		},
		cfg.CommandConfigAction{
			Target: env.MachineName, // will be "" on linux
			Command:"docker run --name vagabond_proxy -d -p 80:80 -v /var/run/docker.sock:/tmp/docker.sock:ro jwilder/nginx-proxy",
		},
		cfg.CommandConfigAction {
			Target: env.MachineName, // will be "" on linux
			Command:fmt.Sprintf("docker run --name vagabond_dnsmasq -d -p 53:53/udp -p 53:53/tcp --cap-add NET_ADMIN andyshinn/dnsmasq --address=/docker/%s", env.MachineIp),
		},
	}
}

func resolverActions(env *config.Environment) ([]cfg.ConfigAction) {
	return []cfg.ConfigAction{
		cfg.ConfigFileAction{
			Filename: "/etc/resolver/docker",
			TemplateVars: env,
			Template: template.Must(template.New("resolver").Parse("nameserver {{.MachineIp}}")),
		},
	}
}