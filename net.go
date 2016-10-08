package main

import (
	"os/exec"
)

func getHostAddress() (string, error) {
	addr, err := exec.Command("defaults", "read", "/Library/Preferences/SystemConfiguration/com.apple.vmnet.plist", "Shared_Net_Address").Output()
	if err != nil {
		return "", err
	}

	return string(addr), nil
}
