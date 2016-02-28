package main

type RebuildCommand struct {
	Disk int `short:"d" long:"disk" description:"size of disk in GiB to create" default:"20"`
}

func (c *RebuildCommand) Execute(args []string) error {
	steps := Steps{
		{
			"Rebuilding disk image",
			func() error {
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
