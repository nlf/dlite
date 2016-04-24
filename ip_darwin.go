package main

import (
	"fmt"
	"os"

	"github.com/nlf/dlite/config"
	"github.com/nlf/dlite/vm"
)

type IPCommand struct{}

func (c *IPCommand) Execute(args []string) error {
	cfg, err := config.New(os.ExpandEnv("$USER"))
	if err != nil {
		return err
	}

	err = cfg.Load()
	if err != nil {
		return err
	}

	v := vm.New(cfg)
	ip, err := v.IP()
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
