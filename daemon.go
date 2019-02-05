package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/kardianos/osext"
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
	d.Error <- fmt.Errorf("shutting down privileged daemon")
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

const template = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
	  <key>Label</key>
		<string>local.dlite</string>
		<key>ProgramArguments</key>
		<array>
		  <string>%s</string>
			<string>daemon</string>
		</array>
		<key>RunAtLoad</key>
		<true/>
  </dict>
</plist>
`

const plistPath = "/Library/LaunchDaemons/local.dlite.plist"

func installDaemon() error {
	exe, err := osext.Executable()
	if err != nil {
		return err
	}

	plist := fmt.Sprintf(template, exe)
	err = ioutil.WriteFile(plistPath, []byte(plist), 0644)
	if err != nil {
		return err
	}

	return exec.Command("launchctl", "load", plistPath).Run()
}

func removeDaemon() error {
	exec.Command("lauunchctl", "unload", plistPath).Run()
	return os.Remove(plistPath)
}
