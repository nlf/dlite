package main

import (
	"github.com/mitchellh/cli"
)

type daemonCommand struct{}

func (c *daemonCommand) Run(args []string) int {
	d := NewDaemon()
	d.Start()
	err := d.Wait()
	if err != nil && err.Error() != "Shutting down privileged daemon" {
		ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *daemonCommand) Synopsis() string {
	return "start the privileged daemon"
}

func (c *daemonCommand) Help() string {
	return "[for internal use] starts the privileged daemon that manages the virtual machine as well as the proxy"
}

func daemonFactory() (cli.Command, error) {
	return &daemonCommand{}, nil
}
