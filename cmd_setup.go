package main

import (
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

		home := ctx.String("home")
		if home == "" {
			return cli.NewExitError("Must specify home", 1)
		}

		err = spin("Setting up DNS", func() error {
			return installResolver(hostname)
		})
		if err != nil {
			return err
		}

		return spin("Setting up NFS", func() error {
			return ensureNFS(home)
		})
	},
}
