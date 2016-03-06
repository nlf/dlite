package main

import "fmt"

type StartCommand struct{}

func (c *StartCommand) Execute(args []string) error {

	steps := Steps{
		{
			"Starting the agent",
			func() error {
				return StartAgent()
			},
		},
	}

	err := Spin(steps)
	if err != nil {
		return err
	}

	fmt.Println("The VM may take some additional time to fully boot")
	return nil
}

func init() {
	var startCommand StartCommand
	cmd.AddCommand("start", "start the daemon", "load and start the launchd agent", &startCommand)
}
