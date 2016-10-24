package main

import (
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

var sshCommand = cli.Command{
	Name:        "ssh",
	Usage:       "start an ssh session with your vm",
	Description: "login to your virtual machine using ssh, this is a convenience shortcut",
	Action: func(ctx *cli.Context) error {
		currentUser := getUser()
		cfg, err := readConfig(getPath(currentUser))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		cmd := exec.Command("ssh", cfg.Hostname)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return cli.NewExitError("", 1)
		}

		return nil
	},
}
