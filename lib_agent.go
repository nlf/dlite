package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/kardianos/osext"
)

const template = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>Label</key>
		<string>local.dlite</string>
		<key>ProgramArguments</key>
		<array>
			<string>/usr/bin/sudo</string>
			<string>%s</string>
			<string>daemon</string>
		</array>
		<key>RunAtLoad</key>
		<true/>
		<key>StandardOutPath</key>
		<string>%s</string>
		<key>StandardErrorPath</key>
		<string>%s</string>
	</dict>
</plist>
`

func CreateAgent() error {
	path, err := osext.Executable()
	if err != nil {
		return err
	}

	fileDir := os.ExpandEnv("$HOME/Library/LaunchAgents")
	err = os.MkdirAll(fileDir, 0755)
	if err != nil {
		return err
	}

	err = changePermissions(fileDir)
	if err != nil {
		return err
	}

	filePath := os.ExpandEnv("$HOME/Library/LaunchAgents/local.dlite.plist")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	outLog := os.ExpandEnv("$HOME/Library/Logs/dlite-out.log")
	errLog := os.ExpandEnv("$HOME/Library/Logs/dlite-err.log")

	plist := fmt.Sprintf(template, path, outLog, errLog)
	_, err = file.WriteString(plist)
	if err != nil {
		return err
	}

	return changePermissions(filePath)
}

func RemoveAgent() error {
	filePath := os.ExpandEnv("$HOME/Library/LaunchAgents/local.dlite.plist")
	return os.RemoveAll(filePath)
}

func StopAgent() error {
	if !AgentRunning() {
		return nil
	}

	filePath := os.ExpandEnv("$HOME/Library/LaunchAgents/local.dlite.plist")
	err := exec.Command("launchctl", "stop", "local.dlite").Run()
	if err != nil {
		return err
	}

	return exec.Command("launchctl", "unload", filePath).Run()
}

func StartAgent() error {
	if AgentRunning() {
		return nil
	}

	filePath := os.ExpandEnv("$HOME/Library/LaunchAgents/local.dlite.plist")
	err := exec.Command("launchctl", "load", filePath).Run()
	if err != nil {
		return err
	}

	return exec.Command("launchctl", "start", "local.dlite").Run()
}

func AgentRunning() bool {
	list, err := exec.Command("sudo", "-u", os.ExpandEnv("$SUDO_USER"), "launchctl", "list", "local.dlite").Output()
	if err != nil {
		return false
	}

	if bytes.Contains(list, []byte("Could not find service")) {
		return false
	}

	pidStartMarker := []byte("\"PID\" = ")
	pidEndMarker := []byte(";")

	pidStart := bytes.Index(list, pidStartMarker)
	if pidStart == -1 {
		return false
	}

	pidEnd := bytes.Index(list[pidStart:], pidEndMarker) + pidStart
	if pidEnd == -1 {
		return false
	}

	pidString := string(list[pidStart+len(pidStartMarker) : pidEnd])
	pid, err := strconv.Atoi(pidString)
	if err != nil {
		return false
	}

	return pid > 0
}
