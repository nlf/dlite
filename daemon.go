package main

type DaemonCommand struct{}

func (c *DaemonCommand) Execute(args []string) error {
	EnsureSudo()
	config, err := ReadConfig()
	if err != nil {
		return err
	}

	err = AddExport(config.Uuid, config.Share)
	if err != nil {
		return err
	}

	StartVM(config)
	ip, err := GetIP(config.Uuid)
	if err != nil {
		return err
	}

	err = AddHost(config.Hostname, ip)
	if err != nil {
		return err
	}

	return Proxy(ip)
}

func init() {
	var daemonCommand DaemonCommand
	cmd.AddCommand("daemon", "internal use", "this command runs the daemon and is intended for internal use only", &daemonCommand)
}
