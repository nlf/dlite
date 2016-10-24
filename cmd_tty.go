package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli"
)

var ttyCommand = cli.Command{
	Name:        "tty",
	Usage:       "open a terminal to the virtual machine",
	Description: "use screen to open a virtual terminal connected to the dlite vm, for most cases SSH is recommended but this can be useful for debugging",
	Action: func(ctx *cli.Context) error {
		currentUser := getUser()

		fmt.Println("You may have to press enter a few times to get a login prompt")
		fmt.Println("Username is 'root', password is 'dlite'")
		fmt.Println("When you're finished press Ctrl+A then D to exit")
		fmt.Println("")
		ttyPath := filepath.Join(getPath(currentUser), "vm.tty")
		cmd := exec.Command("screen", ttyPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := exec.Command("screen", ttyPath).Run()
		if err != nil {
			return cli.NewExitError("Failed to get a terminal, is your virtual machine running?", 1)
		}

		return nil
	},
}
