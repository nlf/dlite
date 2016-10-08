package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Daemon struct {
	Proxy *Proxy
	API   *API
	VM    *VM
	Error chan error
}

func (d *Daemon) Start() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-shutdown
		ui.Info("Got a shutdown signal, halting daemon")
		d.Shutdown()
	}()

	go func() {
		err := d.Proxy.Listen()
		if err != nil {
			fmt.Println("error starting proxy")
			ui.Error(err.Error())
			d.Error <- err
			d.Shutdown()
		}
	}()

	go func() {
		err := d.API.Listen()
		if err != nil {
			ui.Error(err.Error())
			d.Error <- err
			d.Shutdown()
		}
	}()
}

func (d *Daemon) Shutdown() {
	if d.VM != nil {
		d.VM.Stop()
	}
	d.Proxy.Stop()
	d.API.Stop()
	d.Error <- fmt.Errorf("Shutting down privileged daemon")
}

func (d *Daemon) Wait() error {
	return <-d.Error
}

func NewDaemon() *Daemon {
	daemon := &Daemon{}
	proxy := NewProxy(daemon)
	api := NewAPI(daemon)

	daemon.Proxy = proxy
	daemon.API = api
	daemon.Error = make(chan error)
	return daemon
}
