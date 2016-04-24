package rpc

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nlf/dlite/config"
	"github.com/nlf/dlite/nfs"
)

type VM struct {
	server *Server
}

func (v *VM) Start(user string, _ *int) error {
	if v.server.cmd != nil && v.server.cmd.Process != nil {
		return fmt.Errorf("Virtual machine is already running")
	}

	v.server.Status = "The virtual machine is booting, please wait"
	cfg, err := config.New(user)
	if err != nil {
		return err
	}

	err = cfg.Load()
	if err != nil {
		return err
	}

	n := nfs.New(cfg)
	err = n.AddExport()
	if err != nil {
		return err
	}

	err = n.Start()
	if err != nil {
		return err
	}

	self, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}

	c := exec.Command(self, "vm")
	v.server.cmd = c
	return c.Start()
}

func (v *VM) Stop(user string, _ *int) error {
	if v.server.cmd == nil || v.server.cmd.Process == nil {
		return fmt.Errorf("Virtual machine is not running")
	}

	v.server.queue = append(v.server.queue, Message{Command: "shutdown"})
	if v.server.cmd != nil {
		v.server.cmd.Wait()
	}

	v.server.Status = "The virtual machine has not been started"
	v.server.cmd = nil
	v.server.queue = []Message{}
	return nil
}

func (v *VM) SetStatus(status string, _ *int) error {
	if v.server.cmd == nil {
		return fmt.Errorf("Attempted to set status of a VM that has not been started")
	}

	v.server.Status = status
	return nil
}

func (v *VM) GetStatus(_ string, status *string) error {
	if v.server.cmd == nil {
		*status = "The virtual machine has not been started"
	} else {
		*status = v.server.Status
	}
	return nil
}

func (v *VM) GetCommand(_ string, command *Message) error {
	fmt.Printf("%+v\n", v.server.queue)
	if len(v.server.queue) > 0 {
		*command, v.server.queue = v.server.queue[0], v.server.queue[1:]
	} else {
		*command = Message{Command: "wait"}
	}
	return nil
}
