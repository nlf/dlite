package main

import (
	"github.com/nlf/dlite/utils"
)

type UninstallCommand struct{}

func (c *UninstallCommand) Execute(args []string) error {
	utils.EnsureSudo()
	steps := utils.Steps{
		{
			"Removing launchd agent",
			func() error {
				utils.StopAgent()
				utils.RemoveHost()
				return utils.RemoveAgent()
			},
		},
		{
			"Removing files",
			func() error {
				err := utils.RemoveSudoer()
				if err != nil {
					return err
				}

				return utils.RemoveDir()
			},
		},
	}

	return utils.Spin(steps)
}

func init() {
	var uninstallCommand UninstallCommand
	cmd.AddCommand("uninstall", "uninstall dlite", "removes dlite from your system", &uninstallCommand)
}
