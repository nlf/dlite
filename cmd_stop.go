package main

import (
	"github.com/urfave/cli"
)

var stopCommand = cli.Command{
	Name:        "stop",
	Usage:       "stop the virtual machine",
	Description: "stop the virtual machine, exits once the process has ended",
	Action: func(ctx *cli.Context) error {
		return spin("Stopping the virtual machine", func() error {
			return stringRequest("stop")
		})
	},
}
