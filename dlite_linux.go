package main

import (
	"fmt"
	"os"
)

func main() {
	cmd.SubcommandsOptional = true
	_, err := cmd.Parse()
	if err != nil {
		os.Exit(1)
	}

	if cmd.Command.Active == nil {
		fmt.Println("Dlite service daemon")
	}
}
