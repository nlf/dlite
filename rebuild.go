package main

type RebuildCommand struct {
	Disk int `short:"d" long:"disk" description:"size of disk in GiB to create"`
}

func (c *RebuildCommand) Execute(args []string) error {
	config, err := ReadConfig()
	if err != nil {
		return err
	}

	steps := Steps{
		{
			"Rebuilding disk image",
			func() error {
				if c.Disk == 0 {
					c.Disk = config.DiskSize
				} else if c.Disk != config.DiskSize {
					config.DiskSize = c.Disk
					err := SaveConfig(config)
					if err != nil {
						return err
					}
				}
				return CreateDisk(c.Disk)
			},
		},
	}
	return Spin(steps)
}

func init() {
	var rebuildCommand RebuildCommand
	cmd.AddCommand("rebuild", "rebuild your vm", "rebuild the disk for your vm to reset any modifications. this will DESTROY ALL DATA inside your vm.", &rebuildCommand)
}
