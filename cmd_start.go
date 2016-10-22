package main

import (
	"github.com/urfave/cli"
)

var startCommand = cli.Command{
	Name:        "start",
	Usage:       "start the virtual machine",
	Description: "start the virtual machine, exits once booting is complete",
	Action: func(ctx *cli.Context) error {
		return spin("Starting the virtual machine", func() error {
			return stringRequest("start")
		})
	},
}
