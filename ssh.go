package main

import (
	"os"
	"os/exec"
	"syscall"
)

type SSHCommand struct{}

func (c *SSHCommand) Execute(args []string) error {
	config, err := ReadConfig()
	if err != nil {
		return err
	}

	path, err := exec.LookPath("ssh")
	if err != nil {
		return err
	}

	args = append([]string{"", config.Hostname}, args...)
	return syscall.Exec(path, args, os.Environ())
}

func init() {
	var sshCommand SSHCommand
	cmd.AddCommand("ssh", "ssh shortcut", "run an ssh client connected to your vm", &sshCommand)
}
