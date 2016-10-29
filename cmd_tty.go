package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/urfave/cli"
)

var ttyCommand = cli.Command{
	Name:        "tty",
	Usage:       "open a terminal to the virtual machine",
	Description: "use screen to open a virtual terminal connected to the dlite vm, for most cases SSH is recommended but this can be useful for debugging",
	Action: func(ctx *cli.Context) error {
		currentUser := getUser()

		ttyPath := filepath.Join(getPath(currentUser), "vm.tty")
		bin, err := exec.LookPath("screen")
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		err = syscall.Exec(bin, []string{"screen", ttyPath}, os.Environ())
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	},
}
