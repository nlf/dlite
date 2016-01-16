package utils

import (
	"fmt"
	"os"

	"github.com/TheNewNormal/libxhyve"
)

func StartVM(config Config) chan error {
	done := make(chan error)
	ptyCh := make(chan string)
	go func(done chan error) {
		args := []string{
			"-A",
			"-c", fmt.Sprintf("%d", config.CpuCount),
			"-m", fmt.Sprintf("%dG", config.Memory),
			"-s", "0:0,hostbridge",
			"-l", "com1,autopty",
			"-s", "31,lpc",
			"-s", "2:0,virtio-net",
			"-s", "4,virtio-blk," + os.ExpandEnv("$HOME/.dlite/disk.img"),
			"-U", config.Uuid,
			"-f", fmt.Sprintf("kexec,%s,%s,%s", os.ExpandEnv("$HOME/.dlite/bzImage"), os.ExpandEnv("$HOME/.dlite/rootfs.cpio.xz"), "console=ttyS0 hostname=dlite uuid="+config.Uuid),
		}

		err := xhyve.Run(args, ptyCh)
		done <- err
	}(done)

	<-ptyCh
	return done
}
