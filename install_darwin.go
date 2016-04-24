package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/nlf/dlite/config"
	"github.com/nlf/dlite/disk"
	"github.com/nlf/dlite/docker"
	dliteos "github.com/nlf/dlite/os"
	"github.com/nlf/dlite/ssh"
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
- Create a launchd agent in '/Library/LaunchDaemons' used to run the privileged daemon
- Store logs from the daemon in '/Library/Logs'
- Create a launchd agent in '~/Library/LaunchAgents' used to start the privileged daemon under this user's context

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

	var cfg *config.Config
	steps := Steps{
		{
			"Creating configuration",
			func() error {
				cfg, err := config.New(os.ExpandEnv("$USER"))
				if err != nil {
					return err
				}

				if c.DockerVersion == "" {
					latest, err := docker.Latest()
					if err != nil {
						return err
					}
					c.DockerVersion = latest
				}

				if c.Cpus == 0 {
					c.Cpus = runtime.NumCPU()
				}

				cfg.CpuCount = c.Cpus
				cfg.DockerVersion = c.DockerVersion
				cfg.Uuid = uuid.NewV1().String()
				cfg.Memory = c.Memory
				cfg.Hostname = c.Hostname
				cfg.DockerVersion = c.DockerVersion
				cfg.Extra = c.Extra
				cfg.DNSServer = c.DNSServer
				cfg.DiskSize = c.Disk
				cfg.Route = c.Route

				return cfg.Save()
			},
		},
		{
			"Building disk image",
			func() error {
				d := disk.New(cfg)
				d.Detach()

				return d.Create()
			},
		},
		{
			"Downloading OS",
			func() error {
				var vers dliteos.Version
				var err error

				if c.Version == "" {
					vers, err = dliteos.Latest()
				} else {
					vers, err = dliteos.Specific(c.Version)
				}
				if err != nil {
					return err
				}

				return dliteos.Download(cfg.Dir, vers)
			},
		},
		{
			"Generating SSH key",
			func() error {
				s := ssh.New(cfg)
				return s.Generate()
			},
		},
	}

	return Spin(steps)
}

func init() {
	var installCommand InstallCommand
	cmd.AddCommand("install", "install dlite", "creates an empty disk image, downloads the os, saves configuration and creates a launchd agent", &installCommand)
}
