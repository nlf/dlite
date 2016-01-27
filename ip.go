package main

import (
	"fmt"

	"github.com/nlf/dlite/utils"
)

type IPCommand struct{}

func (c *IPCommand) Execute(args []string) error {
	config, err := utils.ReadConfig()
	if err != nil {
		return err
	}

	ip, err := utils.GetIP(config.Uuid)
	if err != nil {
		return err
	}

	fmt.Println(ip)
	return nil
}

func init() {
	var ipCommand IPCommand
	cmd.AddCommand("ip", "get the vm's IP", "display the virtual machine's IP address", &ipCommand)
}
