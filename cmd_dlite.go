package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

const VERSION = "2.0.0-beta7"

// color is overrated but i'm not ready to delete this in case i change my mind
// var ui = &cli.ColoredUi{
// 	InfoColor:  cli.UiColorBlue,
// 	WarnColor:  cli.UiColorYellow,
// 	ErrorColor: cli.UiColorRed,
// 	Ui: &cli.BasicUi{
// 		Writer:      os.Stdout,
// 		ErrorWriter: os.Stderr,
// 		Reader:      os.Stdin,
// 	},
// }

var ui = &cli.BasicUi{
	Writer:      os.Stdout,
	ErrorWriter: os.Stderr,
	Reader:      os.Stdin,
}

func main() {
	app := cli.NewCLI("dlite", VERSION)
	app.Args = os.Args[1:]
	app.Commands = map[string]cli.CommandFactory{
		"daemon": daemonFactory,
		"init":   initFactory,
		"start":  startFactory,
		"stop":   stopFactory,
		"ip":     ipFactory,
		"status": statusFactory,
	}

	status, err := app.Run()
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(status)
}
