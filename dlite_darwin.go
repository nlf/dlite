package main

import (
	"fmt"
	"os"

	"github.com/nlf/dlite/config"
)

func main() {
	cmd.SubcommandsOptional = true
	_, err := cmd.Parse()
	if err != nil {
		os.Exit(1)
	}

	if cmd.Command.Active == nil {
		cfg, err := config.New(os.ExpandEnv("$USER"))
		if err != nil {
			panic(err)
		}

		err = cfg.Load()
		if err != nil {
			panic(err)
		}

		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("DLite has not been installed. Please run 'sudo dlite install'")
			} else {
				fmt.Println(err)
			}
			os.Exit(1)
		}

		fmt.Println("DLite configuration:")
		fmt.Printf("  uuid           : %s\n", cfg.Uuid)
		fmt.Printf("  cpu count      : %d\n", cfg.CpuCount)
		fmt.Printf("  memory         : %d GiB\n", cfg.Memory)
		fmt.Printf("  disk size      : %d GiB\n", cfg.DiskSize)
		fmt.Printf("  hostname       : %s\n", cfg.Hostname)
		fmt.Printf("  dns server     : %s\n", cfg.DNSServer)
		fmt.Printf("  docker version : %s\n", cfg.DockerVersion)
		if cfg.Extra != "" {
			fmt.Printf("  docker args    : %s\n", cfg.Extra)
		}
	}
}
