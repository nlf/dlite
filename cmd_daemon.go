package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mitchellh/cli"
)

type daemonCommand struct{}

func (c *daemonCommand) Run(args []string) int {
	p := NewProxy()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-shutdown
		ui.Info("Got a shutdown signal, halting daemon")
		p.Stop()
	}()

	err := p.Listen()
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *daemonCommand) Synopsis() string {
	return "start the privileged daemon"
}

func (c *daemonCommand) Help() string {
	return "[for internal use] starts the privileged daemon that manages the virtual machine as well as the proxy"
}

func daemonFactory() (cli.Command, error) {
	return &daemonCommand{}, nil
}
