package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func GenerateSSHKey() error {
	path := os.ExpandEnv("$HOME/.dlite/docker")
	if _, err := os.Stat(path); err == nil {
		os.RemoveAll(path)
		os.RemoveAll(path + ".pub")
	}
	output, err := exec.Command("ssh-keygen", "-t", "RSA", "-b", "4096", "-C", "dlite", "-f", path, "-N", "").CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}

	err = changePermissions(path)
	if err != nil {
		return err
	}

	return changePermissions(path + ".pub")
}

func AddSSHConfig(hostname, ip string) error {
	path := os.ExpandEnv("$HOME/.ssh/config")
	key := os.ExpandEnv("$HOME/.dlite/docker")

	config, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if config == nil {
		config = make([]byte, 0)
	}

	startMarker := []byte("# begin dlite")
	endMarker := []byte("# end dlite\n")

	begin := bytes.Index(config, startMarker)
	end := bytes.Index(config, endMarker)

	var temp []byte

	if begin > -1 && end > -1 {
		temp = append(config[:begin], config[end+len(endMarker):]...)
		temp = append(bytes.TrimSpace(temp), '\n')
	} else {
		temp = config
	}

	if len(temp) > 0 && !bytes.HasSuffix(temp, []byte("\n")) {
		temp = append(temp, []byte("\n")...)
	}

	entry := fmt.Sprintf("# begin dlite\nHost %s\n  HostName %s\n  IdentityFile %s\n  User docker\n  StrictHostKeyChecking no\n# end dlite\n", hostname, ip, key)
	temp = append(temp, []byte(entry)...)
	return ioutil.WriteFile(path, temp, 0644)
}

func RemoveSSHConfig() error {
	path := os.ExpandEnv("$HOME/.ssh/config")
	config, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	startMarker := []byte("# begin dlite")
	endMarker := []byte("# end dlite\n")

	begin := bytes.Index(config, startMarker)
	end := bytes.Index(config, endMarker)

	if begin == -1 && end == -1 {
		return nil
	}

	temp := append(config[:begin], config[end+len(endMarker):]...)
	temp = append(bytes.TrimSpace(temp), '\n')
	if len(temp) > 0 && !bytes.HasSuffix(temp, []byte("\n")) {
		temp = append(temp, []byte("\n")...)
	}
	return ioutil.WriteFile(path, temp, 0644)
}

func ShutdownVM() error {
	config, err := ReadConfig()
	if err != nil {
		return err
	}

	return exec.Command("sudo", "-u", os.ExpandEnv("$SUDO_USER"), "ssh", config.Hostname, "sudo", "/sbin/halt").Run()
}
