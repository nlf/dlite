package main

import (
	"github.com/nlf/dlite/utils"
)

type UninstallCommand struct{}

func (c *UninstallCommand) Execute(args []string) error {
	utils.EnsureSudo()
	fmap := utils.FunctionMap{}
	fmap["Removing files"] = func() error {
		err := utils.RemoveSudoer()
		if err != nil {
			return err
		}

		return utils.RemoveDir()
	}

	fmap["Removing launchd agent"] = func() error {
		return utils.RemoveAgent()
	}

	return utils.Spin(fmap)
}

func init() {
	var uninstallCommand UninstallCommand
	cmd.AddCommand("uninstall", "uninstall dlite", "removes dlite from your system", &uninstallCommand)
}
