package utils

import (
	"fmt"
	"os"

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
  </dict>
</plist>
`

func CreateAgent() error {
	path, err := osext.Executable()
	if err != nil {
		return err
	}

	filePath := os.ExpandEnv("$HOME/Library/LaunchAgents/local.dlite.plist")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	plist := fmt.Sprintf(template, path)
	_, err = file.WriteString(plist)
	return err
}

func RemoveAgent() error {
	filePath := os.ExpandEnv("$HOME/Library/LaunchAgents/local.dlite.plist")
	return os.RemoveAll(filePath)
}
