package main

import (
	"os/exec"
	"strings"
)

func getHostAddress() (string, error) {
	addr, err := exec.Command("defaults", "read", "/Library/Preferences/SystemConfiguration/com.apple.vmnet.plist", "Shared_Net_Address").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(addr)), nil
}
