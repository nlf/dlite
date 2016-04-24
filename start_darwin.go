package main

import (
	"os"

	"github.com/nlf/dlite/rpc"
)

type StartCommand struct {
}

func (c *StartCommand) Execute(args []string) error {
	steps := Steps{
		{
			"Starting virtual machine",
			func() error {
				client, err := rpc.NewClient(true)
				if err != nil {
					return err
				}

				var reply int
				return client.Call("VM.Start", os.ExpandEnv("$USER"), &reply)
			},
		},
	}
	return Spin(steps)
}

func init() {
	var startCommand StartCommand
	cmd.AddCommand("start", "start your vm", "start your virtual machine", &startCommand)
}
