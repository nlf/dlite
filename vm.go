package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/TheNewNormal/libxhyve"
	"github.com/tsuru/config"
)

var addressRe = regexp.MustCompile(`.*name=([a-fA-F0-9\-]+)\n.*ip_address=([0-9\.]+)`)

type result struct {
	value string
	err   error
}

type VM struct {
	tty    string
	id     string
	cpu    int
	memory int
	dns    string
	docker string
	extra  string
	uid    int
	gid    int
	sshKey string
	disk   *Disk
	kernel string
	initrd string
}

func (v *VM) Start() (chan error, error) {
	err := v.disk.Attach()
	if err != nil {
		return nil, err
	}

	dev, err := v.disk.RawDevice()
	if err != nil {
		return nil, err
	}

	cmdline := fmt.Sprintf(
		"console=ttyS0 hostname=dlite uuid=%s dns_server=%s ssh_key=%s docker_version=%s docker_extra=%s",
		v.id,
		v.dns,
		v.sshKey,
		v.docker,
		v.extra,
	)

	args := []string{
		"-A",
		"-c", fmt.Sprintf("%d", v.cpu),
		"-m", fmt.Sprintf("%dG", v.memory),
		"-s", "0:0,hostbridge",
		"-l", "com1,autopty",
		"-s", "31,lpc",
		"-s", "2:0,virtio-net",
		"-s", fmt.Sprintf("4,ahci-hd,%s", dev),
		"-U", v.id,
		"-f", fmt.Sprintf(
			"kexec,%s,%s,%s",
			v.kernel,
			v.initrd,
			cmdline,
		),
	}

	pty := make(chan string)
	done := make(chan error)

	go func(args []string, pty chan string, done chan error) {
		done <- xhyve.Run(args, pty)
	}(args, pty, done)

	v.tty = <-pty
	return done, os.Chown(v.tty, v.uid, v.gid)
}

func (v *VM) IP() (*net.TCPAddr, error) {
	value := make(chan result, 1)

	go func() {
		attempts := 0
		for {
			if attempts >= 15 {
				value <- result{"", fmt.Errorf("Timed out waiting for IP address")}
				break
			}

			time.Sleep(time.Second)

			file, err := os.Open("/var/db/dhcpd_leases")
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}

				value <- result{"", err}
				break
			}

			defer file.Close()
			leases, err := ioutil.ReadAll(file)
			if err != nil {
				value <- result{"", err}
				break
			}

			matches := addressRe.FindAllStringSubmatch(string(leases), -1)
			for _, match := range matches {
				if match[1] == v.id {
					value <- result{match[2], nil}
					break
				}
			}

			attempts++
		}
	}()

	res := <-value
	if res.err != nil {
		return nil, res.err
	}

	return net.ResolveTCPAddr("tcp", res.value+":2375")
}

func NewVM() (*VM, error) {
	err := config.ReadConfigFile(filepath.Join(getPath(), "config.yaml"))
	if err != nil {
		return nil, err
	}

	user := getUser()
	path := getPath()

	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return nil, err
	}

	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		return nil, err
	}

	diskPath := filepath.Join(path, "disk.sparseimage")
	diskSize, err := config.GetInt("disk")
	if err != nil {
		return nil, err
	}

	disk, err := NewDisk(diskPath, diskSize, user)
	if err != nil {
		return nil, err
	}

	id, err := config.GetString("id")
	if err != nil {
		return nil, err
	}

	dns, err := config.GetString("dns")
	if err != nil {
		return nil, err
	}

	docker, err := config.GetString("docker")
	if err != nil {
		return nil, err
	}

	extra, err := config.GetString("extra")
	if err != nil {
		return nil, err
	}

	cpu, err := config.GetInt("cpu")
	if err != nil {
		return nil, err
	}

	memory, err := config.GetInt("memory")
	if err != nil {
		return nil, err
	}

	sshBytes, err := ioutil.ReadFile(filepath.Join(path, "key.pub"))
	if err != nil {
		return nil, err
	}
	sshKey := base64.StdEncoding.EncodeToString(sshBytes)

	return &VM{
		id:     id,
		uid:    uid,
		gid:    gid,
		cpu:    cpu,
		memory: memory,
		dns:    dns,
		docker: docker,
		extra:  extra,
		sshKey: sshKey,
		disk:   disk,
		kernel: filepath.Join(path, "bzImage"),
		initrd: filepath.Join(path, "rootfs.cpio.xz"),
	}, nil
}
