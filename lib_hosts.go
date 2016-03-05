package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
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
	path := "/etc/hosts"
	hosts, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	startMarker := []byte("# begin dlite")
	endMarker := []byte("# end dlite\n")

	begin := bytes.Index(hosts, startMarker)
	end := bytes.Index(hosts, endMarker)

	var temp []byte
	if begin > -1 && end > -1 {
		temp = append(hosts[:begin], hosts[end+len(endMarker):]...)
		temp = append(bytes.TrimSpace(temp), '\n')
	} else {
		temp = hosts
	}

	if len(temp) > 0 && !bytes.HasSuffix(temp, []byte("\n")) {
		temp = append(temp, []byte("\n")...)
	}

	entry := fmt.Sprintf("# begin dlite\n%s %s\n# end dlite\n", ip, hostname)
	temp = append(temp, []byte(entry)...)
	return ioutil.WriteFile(path, temp, 0644)
}

func RemoveHost() error {
	path := "/etc/hosts"
	hosts, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	startMarker := []byte("# begin dlite")
	endMarker := []byte("# end dlite\n")

	begin := bytes.Index(hosts, startMarker)
	end := bytes.Index(hosts, endMarker)

	if begin == -1 && end == -1 {
		return nil
	}

	temp := append(hosts[:begin], hosts[end+len(endMarker):]...)
	temp = append(bytes.TrimSpace(temp), '\n')
	if len(temp) > 0 && !bytes.HasSuffix(temp, []byte("\n")) {
		temp = append(temp, []byte("\n")...)
	}
	return ioutil.WriteFile(path, temp, 0644)
}
