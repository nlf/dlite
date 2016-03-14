package main

import (
	"fmt"
	"os"

	"github.com/nlf/libxhyve"
)

func StartVM(config Config) chan error {
	done := make(chan error)
	ptyCh := make(chan string)
	go func(done chan error) {
		path, err := AttachDisk()
		if err != nil {
			done <- err
			return
		}

		user_name := os.ExpandEnv("$SUDO_USER")
		user_id := os.ExpandEnv("$SUDO_UID")
		group_id := os.ExpandEnv("$SUDO_GID")
		home := os.ExpandEnv("$HOME")
		cmdline := fmt.Sprintf("console=ttyS0 hostname=dlite uuid=%s dns_server=%s user_name=%s user_id=%s docker_version=%s docker_extra=%s", config.Uuid, config.DNSServer, user_name, user_id, config.DockerVersion, config.Extra)

		args := []string{
			"-A",
			"-c", fmt.Sprintf("%d", config.CpuCount),
			"-m", fmt.Sprintf("%dG", config.Memory),
			"-s", "0:0,hostbridge",
			"-l", "com1,autopty",
			"-s", "31,lpc",
			"-s", "2:0,virtio-net",
			"-s", "4,ahci-hd," + path,
			"-s", fmt.Sprintf("5,virtio-9p,host=%s,uid=%s,gid=%s", home, user_id, group_id),
			"-U", config.Uuid,
			"-f", fmt.Sprintf("kexec,%s/.dlite/bzImage,%s/.dlite/rootfs.cpio.xz,%s", home, home, cmdline),
		}

		err = xhyve.Run(args, ptyCh)
		done <- err
	}(done)

	socket := <-ptyCh
	changePermissions(socket)
	return done
}
