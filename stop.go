package main

type StopCommand struct{}

func (c *StopCommand) Execute(args []string) error {

	steps := Steps{
		{
			"Stopping the agent",
			func() error {
				return StopAgent()
			},
		},
	}

	return Spin(steps)
}

func init() {
	var stopCommand StopCommand
	cmd.AddCommand("stop", "stop the daemon", "stop and unload the launchd agent", &stopCommand)
}
