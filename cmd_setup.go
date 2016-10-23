package main

import (
	"fmt"

	"github.com/urfave/cli"
)

var setupCommand = cli.Command{
	Name:   "setup",
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

		if err := spin(fmt.Sprintf("Creating /etc/resolver/%s", domain), func() error {
			return installResolver(hostname)
		}); err.ExitCode() != 0 {
			return err
		}

		return spin("Modifying /etc/exports", func() error {
			return ensureNFS(home)
		})
	},
}
