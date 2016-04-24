package vm

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/TheNewNormal/libxhyve"
	"github.com/nlf/dlite/disk"
)

func (v *VM) Start() (chan error, error) {
	d := disk.New(v.config)
	err := d.Attach()
	if err != nil {
		return nil, err
	}

	dev, err := d.RawDevice()
	if err != nil {
		return nil, err
	}

	sshBytes, err := ioutil.ReadFile(filepath.Join(v.config.Dir, "key.pub"))
	if err != nil {
		return nil, err
	}

	sshKey := base64.StdEncoding.EncodeToString(sshBytes)
	cmdline := fmt.Sprintf("console=ttyS0 hostname=dlite uuid=%s dns_server=%s user_name=%s user_id=%d ssh_key=%s docker_version=%s docker_extra=%s",
		v.config.Uuid,
		v.config.DNSServer,
		v.config.Username,
		v.config.Uid,
		sshKey,
		v.config.DockerVersion,
		v.config.Extra,
	)

	args := []string{
		"-A",
		"-c", fmt.Sprintf("%d", v.config.CpuCount),
		"-m", fmt.Sprintf("%dG", v.config.Memory),
		"-s", "0:0,hostbridge",
		"-l", "com1,autopty",
		"-s", "31,lpc",
		"-s", "2:0,virtio-net",
		"-s", fmt.Sprintf("4,ahci-hd,%s", dev),
		"-U", v.config.Uuid,
		"-f", fmt.Sprintf("kexec,%s,%s,%s", filepath.Join(v.config.Dir, "bzImage"), filepath.Join(v.config.Dir, "rootfs.cpio.xz"), cmdline),
	}

	pty := make(chan string)
	done := make(chan error)

	go func(args []string, pty chan string, done chan error) {
		done <- xhyve.Run(args, pty)
	}(args, pty, done)

	v.tty = <-pty
	return done, os.Chown(v.tty, v.config.Uid, v.config.Gid)
}
