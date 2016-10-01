package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/TheNewNormal/libxhyve"
)

var addressRe = regexp.MustCompile(`.*name=([a-fA-F0-9\-]+)\n.*ip_address=([0-9\.]+)`)

type result struct {
	value string
	err   error
}

func startVM() (string, chan error, error) {
	currentUser := getUser()
	configPath := getPath()
	diskPath := filepath.Join(configPath, "disk.sparseimage")

	cfg, err := readConfig()
	if err != nil {
		return "", nil, err
	}

	err = attachDisk(diskPath)
	if err != nil {
		return "", nil, err
	}

	_, dev, err := findDisk(diskPath)
	if err != nil {
		return "", nil, err
	}

	sshBytes, err := ioutil.ReadFile(filepath.Join(configPath, "key.pub"))
	if err != nil {
		return "", nil, err
	}
	sshKey := base64.StdEncoding.EncodeToString(sshBytes)

	cmdline := fmt.Sprintf(
		"console=ttyS0 hostname=dlite uuid=%s dns_server=%s ssh_key=%s docker_version=%s docker_extra=%s",
		cfg.Id,
		cfg.DNS,
		sshKey,
		cfg.Docker,
		cfg.Extra,
	)

	args := []string{
		"-A",
		"-c", fmt.Sprintf("%d", cfg.Cpu),
		"-m", fmt.Sprintf("%dG", cfg.Memory),
		"-s", "0:0,hostbridge",
		"-l", "com1,autopty",
		"-s", "31,lpc",
		"-s", "2:0,virtio-net",
		"-s", fmt.Sprintf("4,ahci-hd,%s", dev),
		"-U", cfg.Id,
		"-f", fmt.Sprintf(
			"kexec,%s,%s,%s",
			filepath.Join(configPath, "bzImage"),
			filepath.Join(configPath, "rootfs.cpio.xz"),
			cmdline,
		),
	}

	pty := make(chan string)
	done := make(chan error)

	go func(args []string, pty chan string, done chan error) {
		done <- xhyve.Run(args, pty)
	}(args, pty, done)

	tty := <-pty
	return tty, done, os.Chown(tty, currentUser.Uid, currentUser.Gid)
}

func stopVM() error {
	return nil
}

func getIP() (*net.TCPAddr, error) {
	cfg, err := readConfig()
	if err != nil {
		return nil, err
	}

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
				if match[1] == cfg.Id {
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
