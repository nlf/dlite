package main

import (
	"os"

	"github.com/urfave/cli"
)

const VERSION = "2.0.0-beta7"

func main() {
	app := cli.NewApp()
	app.Version = VERSION
	app.Usage = "the easiest way to use docker on macOS"
	app.HideHelp = true
	app.UsageText = "dlite <command>"

	app.Commands = []cli.Command{
		daemonCommand,
		setupCommand,
		initCommand,
		startCommand,
		stopCommand,
		statusCommand,
		ipCommand,
		sshCommand,
		ttyCommand,
	}

	app.Run(os.Args)
}
