package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/johanneswuerbach/nfsexports"
)

type result struct {
	value string
	err   error
}

func GetIP(uuid string) (string, error) {
	var ip string
	timeout := make(chan bool, 1)
	value := make(chan result, 1)

	go func() {
		time.Sleep(15 * time.Second)
		timeout <- true
	}()

	go func() {
		for ip == "" {
			time.Sleep(time.Second)
			matchRe := regexp.MustCompile(`.*name=([a-fA-F0-9\-]+)\n.*ip_address=([0-9\.]+)`)

			file, err := os.Open("/var/db/dhcpd_leases")
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				value <- result{"", err}
				return
			}

			defer file.Close()
			leases, err := ioutil.ReadAll(file)
			if err != nil {
				value <- result{"", err}
				return
			}

			lines := string(leases)
			matches := matchRe.FindAllStringSubmatch(lines, -1)
			for _, match := range matches {
				if match[1] == uuid {
					value <- result{match[2], nil}
					break
				}
			}
		}
	}()

	select {
	case res := <-value:
		return res.value, res.err
	case <-timeout:
		return "", fmt.Errorf("Failed to find an IP for the virtual machine in 15 seconds")
	}
}

func AddHost(hostname, ip string) error {
	if hostname == "" {
		hostname = "local.docker"
	}

	ipRe := regexp.MustCompile(`.*# added by dlite$`)
	ipLine := ip + " " + hostname + " # added by dlite"

	file, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	hosts, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	lines := strings.Split(string(hosts), "\n")
	added := false
	for i, line := range lines {
		if ipRe.MatchString(line) {
			lines[i] = ipLine
			added = true
			break
		}
	}

	if !added {
		lines = append(lines, ipLine)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(strings.Join(lines, "\n")))
	return err
}

func RemoveHost() error {
	ipRe := regexp.MustCompile(`.*# added by dlite$`)

	file, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	hosts, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	lines := strings.Split(string(hosts), "\n")
	for i, line := range lines {
		if ipRe.MatchString(line) {
			lines = append(lines[:i], lines[i+1:]...)
			break
		}
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	n, err := file.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return err
	}

	return file.Truncate(int64(n))
}

func AddExport(uuid, share string) error {
	if share == "" {
		share = "/Users"
	}

	export := fmt.Sprintf("%s -network 192.168.64.0 -mask 255.255.255.0 -alldirs -mapall=%s:%s", share, os.Getenv("SUDO_UID"), os.Getenv("SUDO_GID"))
	_, err := nfsexports.Add("", "dlite", export)
	if err != nil {
		return err
	}

	err = nfsexports.ReloadDaemon()
	if err != nil {
		return exec.Command("sudo", "nfsd", "start").Run()
	}

	return err
}
