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
	DNS   *DNS
	Error chan error
}

func (d *Daemon) Start() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-shutdown
		d.Shutdown()
	}()

	go func() {
		err := d.Proxy.Listen()
		if err != nil {
			if err.Error() != "Server closed" {
				d.Error <- err
				d.Shutdown()
			} else {
				d.Error <- nil
			}
		}
	}()

	go func() {
		err := d.API.Listen()
		if err != nil {
			if err.Error() != "Server closed" {
				d.Error <- err
				d.Shutdown()
			} else {
				d.Error <- nil
			}
		}
	}()

	go func() {
		err := d.DNS.Start()
		if err != nil {
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
	d.DNS.Stop()
	d.Error <- fmt.Errorf("Shutting down privileged daemon")
}

func (d *Daemon) Wait() []error {
	err1 := <-d.Error
	err2 := <-d.Error
	return []error{err1, err2}
}

func NewDaemon() *Daemon {
	daemon := &Daemon{}
	proxy := NewProxy(daemon)
	api := NewAPI(daemon)
	dns := NewDNS(daemon)

	daemon.Proxy = proxy
	daemon.API = api
	daemon.DNS = dns
	daemon.Error = make(chan error, 2)
	return daemon
}
