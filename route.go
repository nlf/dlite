package main

import (
	"fmt"
)

type RouteCommand struct{}

func (c *RouteCommand) Execute(args []string) error {
	EnsureSudo()

	if !AgentRunning() {
		return fmt.Errorf("DLite must be running to add routing")
	}

	config, err := ReadConfig()
	if err != nil {
		return err
	}

	return AddRoute(config)
}

func init() {
	var routeCommand RouteCommand
	cmd.AddCommand("route", "enable seamless routing", "configures your network interfaces to allow direct access to containers", &routeCommand)
}
