package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

const VERSION = "2.0.0-beta7"

var ui = &cli.ColoredUi{
	InfoColor:  cli.UiColorBlue,
	WarnColor:  cli.UiColorYellow,
	ErrorColor: cli.UiColorRed,
	Ui: &cli.BasicUi{
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
		Reader:      os.Stdin,
	},
}

func main() {
	app := cli.NewCLI("dlite", VERSION)
	app.Args = os.Args[1:]
	app.Commands = map[string]cli.CommandFactory{
		"daemon": daemonFactory,
		"init":   initFactory,
	}

	status, err := app.Run()
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(status)
}
