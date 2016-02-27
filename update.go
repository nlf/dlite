package main

type UpdateCommand struct {
	Version string `short:"v" long:"version" description:"version of DhyveOS to install"`
}

func (c *UpdateCommand) Execute(args []string) error {
	steps := Steps{
		{
			"Downloading OS",
			func() error {
				if c.Version == "" {
					latest, err := GetLatestOSVersion()
					if err != nil {
						return err
					}
					c.Version = latest
				}
				return DownloadOS(c.Version)
			},
		},
	}

	return Spin(steps)
}

func init() {
	var updateCommand UpdateCommand
	cmd.AddCommand("update", "update your vm", "updates the OS powering your vm", &updateCommand)
}
