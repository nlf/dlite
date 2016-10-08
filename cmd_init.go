package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/satori/go.uuid"
	"github.com/tsuru/config"
)

type initCommand struct{}

func (c *initCommand) Run(args []string) int {
	currentUser := getUser()
	configPath := getPath(currentUser)
	binPath := filepath.Join(configPath, "bin")
	configFile := filepath.Join(configPath, "config.yaml")
	diskFile := filepath.Join(configPath, "disk.qcow")
	cfg := Config{}

	err := config.ReadConfigFile(configFile)
	if err == nil {
		ui.Warn("WARNING: It appears you have already initialized dlite. Continuing will destroy your current virtual machine and its configuration.")
		response, err := ui.Ask("Continue? (y/n)")
		if err != nil {
			ui.Error(err.Error())
			return 1
		}
		response = strings.ToLower(response)

		if response != "y" && response != "yes" {
			ui.Info("Aborting initialization..")
			return 1
		}
	}

	err = os.RemoveAll(configPath)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	err = os.MkdirAll(configPath, 0755)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	cfg.Id = uuid.NewV1().String()

	cfg.Hostname, err = promptString("Virtual machine hostname", "local.docker")
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	cfg.Disk, err = promptInt("Disk size (in gigabytes)", 20)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	cfg.Cpu, err = promptInt("CPU cores to allocate to VM", 2)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	cfg.Memory, err = promptInt("Memory to allocate to VM (in gigabytes)", 2)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	host, err := getHostAddress()
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	cfg.DNS, err = promptString("DNS server", host)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	cfg.Docker, err = promptString("Docker version", "latest")
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	cfg.Extra, err = promptString("Extra flags to pass to the docker daemon", "")
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	fmt.Println("")

	err = spin("Saving configuration", func() error {
		return writeConfig(configPath, cfg)
	})
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	err = spin("Creating tool binaries", func() error {
		err := os.MkdirAll(binPath, 0755)
		if err != nil {
			return err
		}

		for _, tool := range []string{"com.docker.hyperkit", "qcow-tool"} {
			bin, err := Asset(tool)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(filepath.Join(binPath, tool), bin, 0755)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	err = spin("Creating disk", func() error {
		return buildDisk(filepath.Join(binPath, "qcow-tool"), diskFile, cfg.Disk, currentUser.Uid, currentUser.Gid)
	})
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	err = spin("Downloading OS", func() error {
		return downloadOS(configPath)
	})
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *initCommand) Synopsis() string {
	return "initialize your dlite installation"
}

func (c *initCommand) Help() string {
	return "creates a new virtual machine for dlite"
}

func initFactory() (cli.Command, error) {
	return &initCommand{}, nil
}
