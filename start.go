package main

import (
	"github.com/nlf/dlite/utils"
)

type StartCommand struct{}

func (c *StartCommand) Execute(args []string) error {
	return utils.StartAgent()
}

func init() {
	var startCommand StartCommand
	cmd.AddCommand("start", "start the daemon", "load and start the launchd agent", &startCommand)
}
