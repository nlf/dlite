package main

import (
	"github.com/nlf/dlite/utils"
)

type UninstallCommand struct{}

func (c *UninstallCommand) Execute(args []string) error {
	utils.EnsureSudo()
	steps := utils.Steps{
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
		{
			"Removing launchd agent",
			func() error {
				return utils.RemoveAgent()
			},
		},
	}

	return utils.Spin(steps)
}

func init() {
	var uninstallCommand UninstallCommand
	cmd.AddCommand("uninstall", "uninstall dlite", "removes dlite from your system", &uninstallCommand)
}
