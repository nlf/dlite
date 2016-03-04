package main

import (
	"fmt"
	"os"
	"os/exec"
)

type ConfigCommand struct{}

func (c *ConfigCommand) Execute(args []string) error {
	editor := os.ExpandEnv("$EDITOR")
	if editor == "" {
		editor = "vim"
	}

	path := os.ExpandEnv("$HOME/.dlite/config.json")
	editor, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	fmt.Println("Stopping agent...")
	StopAgent()
	fmt.Printf("Editing config file at %s\n", path)
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	fmt.Println("Restarting agent...")
	StartAgent()
	return nil
}

func init() {
	var configCommand ConfigCommand
	cmd.AddCommand("config", "edit config", "edit your vm's configuration", &configCommand)
}
