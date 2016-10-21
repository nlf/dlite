package main

import (
	"net"
	"os/exec"
	"strings"
)

func getNetAddress() (string, error) {
	rawAddr, err := getHostAddress()
	if err != nil {
		return "", err
	}
	addr := net.ParseIP(rawAddr)

	rawMask, err := getNetMask()
	if err != nil {
		return "", err
	}
	mask := net.IPMask(net.ParseIP(rawMask).To4())

	return addr.Mask(mask).String(), nil
}

func getHostAddress() (string, error) {
	addr, err := exec.Command("defaults", "read", "/Library/Preferences/SystemConfiguration/com.apple.vmnet.plist", "Shared_Net_Address").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(addr)), nil
}

func getNetMask() (string, error) {
	mask, err := exec.Command("defaults", "read", "/Library/Preferences/SystemConfiguration/com.apple.vmnet.plist", "Shared_Net_Mask").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(mask)), nil
}
