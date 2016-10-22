package main

import (
	"fmt"

	"github.com/urfave/cli"
)

var statusCommand = cli.Command{
	Name:        "status",
	Usage:       "print the virtual machine's status",
	Description: "fetch and print the configuration and status of the virtual machine",
	Action: func(ctx *cli.Context) error {
		status, err := statusRequest()
		if err != nil {
			return err
		}

		if status.Started {
			fmt.Println("vm_state:       started")
			fmt.Printf("ip_address:     %s\n", status.IP)
			fmt.Printf("pid:            %d\n", status.Pid)
		} else {
			fmt.Println("vm_state:       stopped")
		}
		fmt.Printf("id:             %s\n", status.Id)
		fmt.Printf("hostname:       %s\n", status.Hostname)
		fmt.Printf("disk_size:      %d\n", status.Disk)
		fmt.Printf("disk_path:      %s\n", status.DiskPath)
		fmt.Printf("cpu_cores:      %d\n", status.Cpu)
		fmt.Printf("memory:         %d\n", status.Memory)
		fmt.Printf("dns_server:     %s\n", status.DNS)
		fmt.Printf("docker_version: %s\n", status.Docker)
		fmt.Printf("docker_args:    %s\n", status.Extra)

		return nil
	},
}
