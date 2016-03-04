package main

type UninstallCommand struct{}

func (c *UninstallCommand) Execute(args []string) error {
	EnsureSudo()
	steps := Steps{
		{
			"Removing launchd agent",
			func() error {
				StopAgent()
				RemoveHost()
				RemoveSSHConfig()
				return RemoveAgent()
			},
		},
		{
			"Removing files",
			func() error {
				err := RemoveSudoer()
				if err != nil {
					return err
				}

				return RemoveDir()
			},
		},
	}

	return Spin(steps)
}

func init() {
	var uninstallCommand UninstallCommand
	cmd.AddCommand("uninstall", "uninstall dlite", "removes dlite from your system", &uninstallCommand)
}
