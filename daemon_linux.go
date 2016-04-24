package main

import (
	"log"
	"net"
	"os/exec"
	"time"

	"github.com/nlf/dlite/rpc"
)

type DaemonCommand struct{}

func (c *DaemonCommand) Execute(args []string) error {
	client, err := rpc.NewClient(false)
	if err != nil {
		return err
	}

	log.Println("Polling for commands from server..")
	for {
		time.Sleep(time.Second)
		var cmd rpc.Message
		err := client.Call("VM.GetCommand", "", &cmd)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				continue
			}

			log.Println("Got an error:", err.Error())
			return err
		}

		log.Println("Attempting to process message:", cmd.Command)
		switch cmd.Command {
		case "shutdown":
			log.Println("Got shutdown command, halting virtual machine")
			var i int
			err = client.Call("VM.SetStatus", "The virtual machine is shutting down, please wait", &i)
			err = exec.Command("halt").Run()
			if err != nil {
				log.Println(err)
			}
		case "wait":
			continue
		default:
			log.Printf("Unknown command: %s\n", cmd)
		}
	}

	return nil
}

func init() {
	var daemonCommand DaemonCommand
	cmd.AddCommand("daemon", "start the daemon", "start the rpc client daemon process", &daemonCommand)
}
