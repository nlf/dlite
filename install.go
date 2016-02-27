package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/satori/go.uuid"
)

type InstallCommand struct {
	Cpus     int    `short:"c" long:"cpus" description:"number of CPUs to allocate" default:"1"`
	Disk     int    `short:"d" long:"disk" description:"size of disk in GiB to create" default:"20"`
	Memory   int    `short:"m" long:"memory" description:"amount of memory in GiB to allocate" default:"2"`
	SSHKey   string `short:"s" long:"ssh-key" description:"path to public ssh key" default:"$HOME/.ssh/id_rsa.pub"`
	Version  string `short:"v" long:"os-version" description:"version of DhyveOS to install"`
	Hostname string `short:"n" long:"hostname" description:"hostname to use for vm" default:"local.docker"`
	Share    string `short:"S" long:"share" description:"directory to export from NFS" default:"/Users"`
}

func (c *InstallCommand) Execute(args []string) error {
	EnsureSudo()

	fmt.Println("The install command will make the following changes to your system:")
	fmt.Println("- Create a '.dlite' directory in your home")
	fmt.Printf("- Create a %dGB disk image in the '.dlite' directory\n", c.Disk)
	if c.Version == "" {
		fmt.Println("- Download the latest version of DhyveOS to the '.dlite' directory")
	} else {
		fmt.Printf("- Download version %s of DhyveOS to the '.dlite' directory\n", c.Version)
	}
	fmt.Println("- Create a 'config.json' file in the '.dlite' directory")
	fmt.Println("- Add a line to your sudoers file to allow running the 'dlite' binary without a password")
	fmt.Println("- Create a launchd agent in '~/Library/LaunchAgents' used to run the daemon")
	fmt.Println("- Store logs from the daemon in '~/Library/Logs'")

	fmt.Print("Would you like to continue? (Y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "n" || response == "no" {
		return fmt.Errorf("Aborted install due to user input")
	}

	if response != "" && response != "y" && response != "yes" {
		return fmt.Errorf("Aborted install due to invalid user input")
	}

	steps := Steps{
		{
			"Building disk image",
			func() error {
				// clean up but ignore errors since it's possible things weren't installed
				StopAgent()
				RemoveAgent()
				RemoveHost()
				RemoveDir()

				err := CreateDir()
				if err != nil {
					return err
				}

				return CreateDisk(c.SSHKey, c.Disk)
			},
		},
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
		{
			"Writing configuration",
			func() error {
				uuid := uuid.NewV1().String()
				return SaveConfig(Config{
					Uuid:     uuid,
					CpuCount: c.Cpus,
					Memory:   c.Memory,
					Hostname: c.Hostname,
					Share:    c.Share,
				})
			},
		},
		{
			"Creating launchd agent",
			func() error {
				err := AddSudoer()
				if err != nil {
					return err
				}

				return CreateAgent()
			},
		},
	}

	return Spin(steps)
}

func init() {
	var installCommand InstallCommand
	cmd.AddCommand("install", "install dlite", "creates an empty disk image, downloads the os, saves configuration and creates a launchd agent", &installCommand)
}
