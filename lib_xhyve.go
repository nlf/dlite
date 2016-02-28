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

		args := []string{
			"-A",
			"-c", fmt.Sprintf("%d", config.CpuCount),
			"-m", fmt.Sprintf("%dG", config.Memory),
			"-s", "0:0,hostbridge",
			"-l", "com1,autopty",
			"-s", "31,lpc",
			"-s", "2:0,virtio-net",
			"-s", "4,ahci-hd," + path,
			"-U", config.Uuid,
			"-f", fmt.Sprintf("kexec,%s,%s,%s", os.ExpandEnv("$HOME/.dlite/bzImage"), os.ExpandEnv("$HOME/.dlite/rootfs.cpio.xz"), "console=ttyS0 hostname=dlite uuid="+config.Uuid+" share="+config.Share),
		}

		err = xhyve.Run(args, ptyCh)
		done <- err
	}(done)

	<-ptyCh
	return done
}
