package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var uninstallCommand = cli.Command{
	Name:        "uninstall",
	Usage:       "completely remove dlite",
	Description: "remove the user's virtual machine as well as the system service and reset all configuration changes",
	Action: func(ctx *cli.Context) error {
		var hostname string
		var home string

		if err := spin("Removing virtual machine", func() error {
			stringRequest("stop")

			currentUser := getUser()
			home = currentUser.Home

			basePath := getPath(currentUser)
			cfg, err := readConfig(basePath)
			if err != nil {
				return err
			}

			hostname = cfg.Hostname

			err = removeSSHConfig(currentUser, hostname)
			if err != nil {
				return err
			}

			return os.RemoveAll(basePath)
		}); err.ExitCode() != 0 {
			return err
		}

		fmt.Println("")
		fmt.Println("Next we'll run a few steps that require sudo, you may be prompted for your password.")
		return runCleanup(hostname, home)
	},
}
