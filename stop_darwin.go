package main

import (
	"os"

	"github.com/nlf/dlite/rpc"
)

type StopCommand struct {
}

func (c *StopCommand) Execute(args []string) error {
	steps := Steps{
		{
			"Stopping virtual machine",
			func() error {
				client, err := rpc.NewClient(true)
				if err != nil {
					return err
				}

				var reply int
				return client.Call("VM.Stop", os.ExpandEnv("$USER"), &reply)
			},
		},
	}
	return Spin(steps)
}

func init() {
	var stopCommand StopCommand
	cmd.AddCommand("stop", "stop your vm", "stop your virtual machine", &stopCommand)
}
