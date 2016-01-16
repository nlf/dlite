package main

import (
	"fmt"
)

type VersionCommand struct{}

func (c *VersionCommand) Execute(args []string) error {
	fmt.Println(VERSION)
	return nil
}

func init() {
	var versionCommand VersionCommand
	cmd.AddCommand("version", "display current version", "displays the currently installed version of dlite", &versionCommand)
}
