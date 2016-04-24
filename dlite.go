package main

import (
	"github.com/jessevdk/go-flags"
)

const VERSION = "2.0.0-beta6"

type Options struct{}

var (
	options Options
	cmd     = flags.NewParser(&options, flags.Default)
)
