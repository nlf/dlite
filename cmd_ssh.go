package main

import (
	"os"
	"os/exec"
	"syscall"

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

		bin, err := exec.LookPath("ssh")
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		err = syscall.Exec(bin, []string{"ssh", cfg.Hostname}, os.Environ())
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	},
}
