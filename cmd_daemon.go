package main

import (
	"github.com/urfave/cli"
)

var daemonCommand = cli.Command{
	Name:   "daemon",
	Hidden: true,
	Action: func(ctx *cli.Context) error {
		d := NewDaemon()
		d.Start()
		errs := d.Wait()
		for _, err := range errs {
			if err != nil && err.Error() != "Shutting down privileged daemon" {
				return cli.NewExitError(err.Error(), 1)
			}
		}

		return nil
	},
}
