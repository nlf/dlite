package main

import (
	"github.com/nlf/dlite/utils"
)

type RebuildCommand struct {
	Disk   int    `short:"d" long:"disk" description:"size of disk in GiB to create" default:"20"`
	SSHKey string `short:"s" long:"ssh-key" description:"path to public ssh key" default:"$HOME/.ssh/id_rsa.pub"`
}

func (c *RebuildCommand) Execute(args []string) error {
	steps := utils.Steps{
		{
			"Rebuilding disk image",
			func() error {
				return utils.CreateDisk(c.SSHKey, c.Disk)
			},
		},
	}
	return utils.Spin(steps)
}

func init() {
	var rebuildCommand RebuildCommand
	cmd.AddCommand("rebuild", "rebuild your vm", "rebuild the disk for your vm to reset any modifications. this will DESTROY ALL DATA inside your vm.", &rebuildCommand)
}
