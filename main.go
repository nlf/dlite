package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

const VERSION = "1.1.1"

type Options struct{}

var (
	options Options
	cmd     = flags.NewParser(&options, flags.Default)
)

func main() {
	cmd.SubcommandsOptional = true
	_, err := cmd.Parse()
	if err != nil {
		os.Exit(1)
	}

	if cmd.Command.Active == nil {
		cmd.WriteHelp(os.Stdout)
		os.Exit(1)
	}
}
