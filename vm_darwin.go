package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nlf/dlite/config"
	"github.com/nlf/dlite/rpc"
	"github.com/nlf/dlite/vm"
)

type VMCommand struct {
	Hidden bool
}

func (c *VMCommand) Execute(args []string) error {
	client, err := rpc.NewClient(true)
	if err != nil {
		return err
	}

	cfg, err := config.New(os.ExpandEnv("$SUDO_USER"))
	if err != nil {
		return err
	}

	err = cfg.Load()
	if err != nil {
		return err
	}

	vm := vm.New(cfg)
	done, err := vm.Start()
	if err != nil {
		return err
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-shutdown
		var i int
		client.Call("VM.Stop", "", &i)
	}()

	err = <-done
	return err
}

func init() {
	var vmCommand VMCommand
	com, _ := cmd.AddCommand("vm", "start a vm", "start the actual vm", &vmCommand)
	com.Hidden = true
}
