package main

import (
	"os"

	"github.com/nlf/dlite/config"
	"github.com/nlf/dlite/disk"
)

type RebuildCommand struct {
	Disk int `short:"d" long:"disk" description:"size of disk in GiB to create"`
}

func (c *RebuildCommand) Execute(args []string) error {
	var cfg *config.Config

	steps := Steps{
		{
			"Reading configuration",
			func() error {
				cfg, err := config.New(os.ExpandEnv("$USER"))
				if err != nil {
					return err
				}

				return cfg.Load()
			},
		},
		{
			"Rebuilding disk image",
			func() error {
				if cfg.DiskSize != c.Disk {
					cfg.DiskSize = c.Disk
					err := cfg.Save()
					if err != nil {
						return err
					}
				}

				d := disk.New(cfg)
				d.Detach()
				return d.Create()
			},
		},
	}
	return Spin(steps)
}

func init() {
	var rebuildCommand RebuildCommand
	cmd.AddCommand("rebuild", "rebuild your vm", "rebuild the disk for your vm to reset any modifications. this will DESTROY ALL DATA inside your vm.", &rebuildCommand)
}
