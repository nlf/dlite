package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func generateKeys(user User) error {
	base := getPath(user)
	output, err := exec.Command("ssh-keygen", "-f", fmt.Sprintf("%s/key", base), "-P", "").CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(output))
	}

	return nil
}

func addSSHConfig(user User, hostname string) error {
	configPath := filepath.Join(user.Home, ".ssh", "config")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := ioutil.WriteFile(configPath, []byte(""), 0644)
		if err != nil {
			return err
		}
	}

	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	hasHost := strings.Index(string(config), fmt.Sprintf("Host %s", hostname))
	if hasHost != -1 {
		return nil
	}

	keyfile := filepath.Join(getPath(user), "key")
	newConfig := string(config)
	newConfig += fmt.Sprintf("Host %s\n  User docker\n  IdentityFile %s", hostname, keyfile)
	return ioutil.WriteFile(configPath, []byte(newConfig), 0644)
}

func removeSSHConfig(user User, hostname string) error {
	configPath := filepath.Join(user.Home, ".ssh", "config")
	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	keyfile := filepath.Join(getPath(user), "key")
	hostConfig := fmt.Sprintf("Host %s\n  User docker\n  IdentityFile %s", hostname, keyfile)

	hostMatcher := regexp.MustCompile(fmt.Sprintf("(?m)^%s?$", hostConfig))
	newConfig := hostMatcher.ReplaceAllString(string(config), "")

	return ioutil.WriteFile(configPath, []byte(newConfig), 0644)
}
