package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nlf/dlite/proxy"
	"github.com/nlf/dlite/rpc"
)

type DaemonCommand struct{}

func (c *DaemonCommand) Execute(args []string) error {
	log.Println("Starting rpc server..")
	server, err := rpc.NewServer()
	if err != nil {
		return err
	}

	prox := proxy.New(server)
	go server.Listen()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-shutdown
		log.Println("Got a shutdown signal, stopping rpc server")
		server.Stop()
		log.Println("Stopping proxy")
		prox.Stop()
	}()

	log.Println("Starting proxy..")
	return prox.Listen()
}

func init() {
	var daemonCommand DaemonCommand
	cmd.AddCommand("daemon", "start the daemon", "start the privileged daemon process", &daemonCommand)
}
