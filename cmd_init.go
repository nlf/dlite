package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/satori/go.uuid"
	"github.com/tsuru/config"
)

type initCommand struct{}

func (c *initCommand) Run(args []string) int {
	currentUser := getUser()
	configPath := getPath()

	err := config.ReadConfigFile(configPath + "/config.yaml")
	if err == nil {
		ui.Warn("WARNING: It appears you have already initialized dlite. Continuing will destroy your current virtual machine.")
		response, err := ui.Ask("Continue? (y/n)")
		if err != nil {
			ui.Error(err.Error())
			return 1
		}
		response = strings.ToLower(response)

		if response != "y" && response != "yes" {
			ui.Info("Aborting initialization..")
			return 1
		}
	}

	err = os.RemoveAll(configPath)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	err = os.MkdirAll(configPath, 0755)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	config.Set("id", uuid.NewV1().String())

	hostname, err := ui.Ask("Virtual machine hostname: [local.docker]")
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	if hostname == "" {
		hostname = "local.docker"
	}

	config.Set("hostname", hostname)

	diskStr, err := ui.Ask("Disk size (in gigabytes): [20]")
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	disk := 20
	if diskStr != "" {
		disk, err = strconv.Atoi(diskStr)
		if err != nil {
			ui.Error(err.Error())
			return 1
		}
	}
	config.Set("disk", disk)

	cpuStr, err := ui.Ask("CPU cores to allocate to VM: [2]")
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	cpu := 1
	if cpuStr != "" {
		cpu, err = strconv.Atoi(cpuStr)
		if err != nil {
			ui.Error(err.Error())
			return 1
		}
	}
	config.Set("cpu", cpu)

	memStr, err := ui.Ask("Memory to allocate to VM (in gigabytes): [2]")
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	mem := 2
	if memStr != "" {
		mem, err = strconv.Atoi(memStr)
		if err != nil {
			ui.Error(err.Error())
			return 1
		}
	}
	config.Set("memory", mem)

	dns, err := ui.Ask("DNS server: [192.168.64.1]")
	if err != nil {
		ui.Error(err.Error())
		return 1
	}
	if dns == "" {
		dns = "192.168.64.1"
	}
	config.Set("dns", dns)

	docker, err := ui.Ask("Requested docker version: [latest]")
	if err != nil {
		ui.Error(err.Error())
		return 1
	}
	if docker == "" {
		docker = "latest"
	}
	config.Set("docker", docker)

	extra, err := ui.Ask("Extra flags to pass to the docker daemon:")
	if err != nil {
		ui.Error(err.Error())
		return 1
	}
	config.Set("extra", extra)

	ui.Info("Saving configuration..")
	err = config.WriteConfigFile(configPath+"/config.yaml", 0644)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	ui.Info("Creating disk..")
	d, err := NewDisk(configPath+"/disk.sparseimage", disk, currentUser)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	err = d.Build()
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	ui.Info("Downloading OS..")
	err = DownloadOS(configPath)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *initCommand) Synopsis() string {
	return "initialize your dlite installation"
}

func (c *initCommand) Help() string {
	return "creates a new virtual machine for dlite"
}

func initFactory() (cli.Command, error) {
	return &initCommand{}, nil
}
