package main

import (
	"github.com/nlf/dlite/rpc"
)

type StatusCommand struct{}

func (c *StatusCommand) Execute(args []string) error {
	client, err := rpc.NewClient(false)
	if err != nil {
		return err
	}

	status := args[0]
	var i int
	return client.Call("VM.SetStatus", status, &i)
}

func init() {
	var statusCommand StatusCommand
	cmd.AddCommand("status", "set the vm status", "use the rpc channel to update the virtual machine's status", &statusCommand)
}
