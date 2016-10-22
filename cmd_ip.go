package main

import (
	"fmt"

	"github.com/urfave/cli"
)

var ipCommand = cli.Command{
	Name:        "ip",
	Usage:       "display the virtual machine's IP",
	Description: "lookup and print the IP address of the virtual machine",
	Action: func(ctx *cli.Context) error {
		status, err := statusRequest()
		if err != nil {
			return err
		}

		if !status.Started {
			return cli.NewExitError("Virtual machine not running", 1)
		}

		fmt.Println(status.IP)
		return nil
	},
}
