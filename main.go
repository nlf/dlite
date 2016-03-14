package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

const VERSION = "2.0.0-beta4"

type Options struct{}

var (
	options Options
	cmd     = flags.NewParser(&options, flags.Default)
)

func main() {
	cmd.SubcommandsOptional = true
	_, err := cmd.Parse()
	if err != nil {
		os.Exit(1)
	}

	if cmd.Command.Active == nil {
		config, err := ReadConfig()
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("DLite has not been installed. Please run 'sudo dlite install'")
			} else {
				fmt.Println(err)
			}
			os.Exit(1)
		}

		fmt.Println("DLite configuration:")
		fmt.Printf("  uuid           : %s\n", config.Uuid)
		fmt.Printf("  cpu count      : %d\n", config.CpuCount)
		fmt.Printf("  memory         : %d GiB\n", config.Memory)
		fmt.Printf("  disk size      : %d GiB\n", config.DiskSize)
		fmt.Printf("  hostname       : %s\n", config.Hostname)
		fmt.Printf("  dns server     : %s\n", config.DNSServer)
		fmt.Printf("  docker version : %s\n", config.DockerVersion)
		if config.Extra != "" {
			fmt.Printf("  docker args    : %s\n", config.Extra)
		}
	}
}
