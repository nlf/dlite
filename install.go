package main

import (
	"github.com/nlf/dlite/utils"
	"github.com/satori/go.uuid"
)

type InstallCommand struct {
	Cpus    int    `short:"c" long:"cpus" description:"number of CPUs to allocate" default:"1"`
	Disk    int    `short:"d" long:"disk" description:"size of disk in GiB to create" default:"20"`
	Memory  int    `short:"m" long:"memory" description:"amount of memory in GiB to allocate" default:"2"`
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

	steps := utils.Steps{
		{
			"Building disk image",
			func() error {
				// clean up but ignore errors since it's possible things weren't installed
				utils.StopAgent()
				utils.RemoveAgent()
				utils.RemoveHost()
				utils.RemoveDir()

				err := utils.CreateDir()
				if err != nil {
					return err
				}

				return utils.CreateDisk(c.SSHKey, c.Disk)
			},
		},
		{
			"Downloading OS",
			func() error {
				if c.Version == "" {
					latest, err := utils.GetLatestOSVersion()
					if err != nil {
						return err
					}
					c.Version = latest
				}
				return utils.DownloadOS(c.Version)
			},
		},
		{
			"Writing configuration",
			func() error {
				uuid := uuid.NewV1().String()
				return utils.SaveConfig(utils.Config{
					Uuid:     uuid,
					CpuCount: c.Cpus,
					Memory:   c.Memory,
					Hostname: c.Hostname,
				})
			},
		},
		{
			"Creating launchd agent",
			func() error {
				err := utils.AddSudoer()
				if err != nil {
					return err
				}

				return utils.CreateAgent()
			},
		},
	}

	return utils.Spin(steps)
}

func init() {
	var installCommand InstallCommand
	cmd.AddCommand("install", "install dlite", "creates an empty disk image, downloads the os, saves configuration and creates a launchd agent", &installCommand)
}
