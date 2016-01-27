package main

import (
	"github.com/nlf/dlite/utils"
)

type StopCommand struct{}

func (c *StopCommand) Execute(args []string) error {
	return utils.StopAgent()
}

func init() {
	var stopCommand StopCommand
	cmd.AddCommand("stop", "stop the daemon", "stop and unload the launchd agent", &stopCommand)
}
