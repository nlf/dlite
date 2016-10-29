package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var cleanupCommand = cli.Command{
	Name:   "cleanup",
	Hidden: true,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "hostname",
		},
		cli.StringFlag{
			Name: "home",
		},
	},
	Action: func(ctx *cli.Context) error {
		err := ensureRoot()
		if err != nil {
			return err
		}

		hostname := ctx.String("hostname")
		if hostname == "" {
			return cli.NewExitError("Must specify hostname", 1)
		}
		domain := getDomain(hostname)

		home := ctx.String("home")
		if home == "" {
			return cli.NewExitError("Must specify home", 1)
		}

		if err := spin(fmt.Sprintf("Removing /etc/resolver/%s", domain), func() error {
			return os.Remove(fmt.Sprintf("/etc/resolver/%s", domain))
		}); err.ExitCode() != 0 {
			return err
		}

		if err := spin("Modifying /etc/exports", func() error {
			return removeNFS(home)
		}); err.ExitCode() != 0 {
			return err
		}

		return spin("Removing /Library/LaunchDaemons/local.dlite.plist", func() error {
			return removeDaemon()
		})
	},
}
