package main

import (
	"os"
	"os/signal"
	"syscall"
)

type DaemonCommand struct{}

func (c *DaemonCommand) Execute(args []string) error {
	EnsureSudo()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, os.Kill)
	go func() {
		<-shutdown
		ShutdownVM()
	}()

	config, err := ReadConfig()
	if err != nil {
		return err
	}

	done := StartVM(config)
	ip, err := GetIP(config.Uuid)
	if err != nil {
		return err
	}

	err = AddHost(config.Hostname, ip)
	if err != nil {
		return err
	}

	err = AddSSHConfig(config.Hostname, ip)
	if err != nil {
		return err
	}

	if config.Route {
		err := AddRoute(config)
		if err != nil {
			return err
		}
	}

	return Proxy(ip, done)
}

func init() {
	var daemonCommand DaemonCommand
	cmd.AddCommand("daemon", "internal use", "this command runs the daemon and is intended for internal use only", &daemonCommand)
}
