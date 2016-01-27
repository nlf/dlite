package main

import (
	"github.com/nlf/dlite/utils"
)

type DaemonCommand struct{}

func (c *DaemonCommand) Execute(args []string) error {
	utils.EnsureSudo()
	config, err := utils.ReadConfig()
	if err != nil {
		return err
	}

	err = utils.AddExport(config.Uuid, config.Share)
	if err != nil {
		return err
	}

	utils.StartVM(config)
	ip, err := utils.GetIP(config.Uuid)
	if err != nil {
		return err
	}

	err = utils.AddHost(config.Hostname, ip)
	if err != nil {
		return err
	}

	return utils.Proxy(ip)
}

func init() {
	var daemonCommand DaemonCommand
	cmd.AddCommand("daemon", "internal use", "this command runs the daemon and is intended for internal use only", &daemonCommand)
}
