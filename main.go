package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var PREFIX string

func main() {
	PREFIX = os.Getenv("PREFIX")
	if PREFIX == "" {
		PREFIX = "/usr/local"
	}

	app := cli.NewApp()
	app.Name = "dlite"
	app.Usage = "The simplest way to run Docker on OS X"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "install dlite",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "cpu, c",
					Value: 1,
					Usage: "number of CPUs to allocate to the vm",
				},
				cli.IntFlag{
					Name:  "disk, d",
					Value: 30,
					Usage: "size of disk (in gigabytes) to create",
				},
				cli.IntFlag{
					Name:  "memory, m",
					Value: 1,
					Usage: "amount of memory (in gigabytes) to allocate to the vm",
				},
				cli.StringFlag{
					Name:  "ssh-key, s",
					Value: "$HOME/.ssh/id_rsa.pub",
					Usage: "path to the public ssh key to add to the vm",
				},
			},
			Action: installHandler,
		},
		{
			Name:   "daemon",
			Usage:  "run the daemon",
			Action: daemonHandler,
		},
	}

	app.Run(os.Args)
}
