package main

import (
	"net"
	"os/exec"
	"strings"
)

func getNetAddress() (string, error) {
	rawAddr, _ := getHostAddress()
	addr := net.ParseIP(rawAddr)

	rawMask, _ := getNetMask()
	mask := net.IPMask(net.ParseIP(rawMask).To4())

	return addr.Mask(mask).String(), nil
}

func getHostAddress() (string, error) {
	addr, err := exec.Command("defaults", "read", "/Library/Preferences/SystemConfiguration/com.apple.vmnet.plist", "Shared_Net_Address").Output()
	if err != nil {
		return "192.168.64.1", err
	}

	return strings.TrimSpace(string(addr)), nil
}

func getNetMask() (string, error) {
	mask, err := exec.Command("defaults", "read", "/Library/Preferences/SystemConfiguration/com.apple.vmnet.plist", "Shared_Net_Mask").Output()
	if err != nil {
		return "255.255.255.0", err
	}

	return strings.TrimSpace(string(mask)), nil
}

func getDomain(hostname string) string {
	var domain string

	lastDot := strings.LastIndex(hostname, ".")
	if lastDot > -1 {
		domain = hostname[lastDot+1:]
	} else {
		domain = hostname
	}

	return domain
}
