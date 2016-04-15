package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/satori/go.uuid"
)

type InstallCommand struct {
	Cpus          int    `short:"c" long:"cpus" description:"number of CPUs to allocate" default-mask:"# of CPUs"`
	Disk          int    `short:"d" long:"disk" description:"size of disk in GiB to create" default:"20"`
	DNSServer     string `short:"s" long:"dns-server" description:"DNS server to use in the vm" default:"192.168.64.1"`
	DockerVersion string `short:"D" long:"docker-version" description:"version of Docker to install"`
	Extra         string `short:"e" long:"extra" description:"extra arguments to pass to Docker"`
	Hostname      string `short:"n" long:"hostname" description:"hostname to use for vm" default:"local.docker"`
	Memory        int    `short:"m" long:"memory" description:"amount of memory in GiB to allocate" default:"2"`
	Version       string `short:"v" long:"os-version" description:"version of DhyveOS to install"`
	Route         bool   `short:"r" long:"route" description:"add routing entries to allow direct connections to containers"`
}

func (c *InstallCommand) Execute(args []string) error {
	EnsureSudo()

	versionMsg := "the latest version"
	if c.Version != "" {
		versionMsg = "version " + c.Version
	}
	fmt.Printf(`
The install command will make the following changes to your system:
- Create a '.dlite' directory in your home
- Create a %d GiB sparse disk image in the '.dlite' directory
- Download %s of DhyveOS to the '.dlite' directory
- Create a 'config.json' file in the '.dlite' directory
- Create a new SSH key pair in the '.dlite' directory for the vm
- Add a line to your sudoers file to allow running the 'dlite' binary without a password
- Create a launchd agent in '~/Library/LaunchAgents' used to run the daemon
- Store logs from the daemon in '~/Library/Logs'

IMPORTANT: if the dlite binary is in a path writeable by the user (which is the case if you installed with Homebrew or go get, for example),
the sudoers change will let any attacker bypass the sudo password by modifying the binary. THIS EFFECTIVELY DISABLES SUDO PASSWORD SECURITY.

In addition to the above actions that take place during installation, when the service is started a few other files are modified.
While DLite makes every effort to not damage any of these files, it is advisable for you to back them up manually before installation
The files are:
- /etc/hosts
- /etc/sudoers
- ~/.ssh/config

`, c.Disk, versionMsg)

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

				return CreateDisk(c.Disk)
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
			"Generating SSH key",
			func() error {
				return GenerateSSHKey()
			},
		},
		{
			"Writing configuration",
			func() error {
				if c.DockerVersion == "" {
					latest, err := GetLatestDockerVersion()
					if err != nil {
						return err
					}
					c.DockerVersion = latest
				}
				if c.Cpus == 0 {
					c.Cpus = runtime.NumCPU()
				}
				uuid := uuid.NewV1().String()
				return SaveConfig(Config{
					Uuid:          uuid,
					CpuCount:      c.Cpus,
					Memory:        c.Memory,
					Hostname:      c.Hostname,
					DockerVersion: c.DockerVersion,
					Extra:         c.Extra,
					DNSServer:     c.DNSServer,
					DiskSize:      c.Disk,
					Route:         c.Route,
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
