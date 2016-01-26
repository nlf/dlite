package main

import (
	"github.com/nlf/dlite/utils"
	"github.com/satori/go.uuid"
)

type InstallCommand struct {
	Cpus    int    `short:"c" long:"cpus" description:"number of CPUs to allocate" default:"1"`
	Disk    int    `short:"d" long:"disk" description:"size of disk in GiB to create" default:"30"`
	Memory  int    `short:"m" long:"memory" description:"amount of memory in GiB to allocate" default:"1"`
	SSHKey  string `short:"s" long:"ssh-key" description:"path to public ssh key" default:"$HOME/.ssh/id_rsa.pub"`
	Version string `short:"v" long:"os-version" description:"version of DhyveOS to install"`
	Hostname string `short:"n" long:"hostname" description:"hostname to use for vm" default:"local.docker"`
}

func (c *InstallCommand) Execute(args []string) error {
	utils.EnsureSudo()
	err := utils.CreateDir()
	if err != nil {
		return err
	}

	fmap := utils.FunctionMap{}
	fmap["Building disk image"] = func() error {
		return utils.CreateDisk(c.SSHKey, c.Disk)
	}

	fmap["Downloading OS"] = func() error {
		if c.Version == "" {
			latest, err := utils.GetLatestOSVersion()
			if err != nil {
				return err
			}
			c.Version = latest
		}
		return utils.DownloadOS(c.Version)
	}

	fmap["Writing configuration"] = func() error {
		uuid := uuid.NewV1().String()
		return utils.SaveConfig(utils.Config{
			Uuid: uuid,
			CpuCount: c.Cpus,
			Memory: c.Memory,
			Hostname: c.Hostname,
		})
	}

	fmap["Creating launchd agent"] = func() error {
		err := utils.AddSudoer()
		if err != nil {
			return err
		}

		return utils.CreateAgent()
	}

	return utils.Spin(fmap)
}

func init() {
	var installCommand InstallCommand
	cmd.AddCommand("install", "install dlite", "creates an empty disk image, downloads the os, saves configuration and creates a launchd agent", &installCommand)
}
