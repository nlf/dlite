package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/DHowett/go-plist"
)

func buildDisk(path string, size, uid, gid int) error {
	err := exec.Command("hdiutil", "create", "-size", fmt.Sprintf("%dg", size), "-type", "SPARSE", "-layout", "MBRSPUD", path).Run()
	if err != nil {
		return err
	}

	return os.Chown(path, uid, gid)
}

func findDisk(path string) (string, string, error) {
	type device struct {
		DevEntry string `plist:"dev-entry"`
	}

	type image struct {
		ImagePath string   `plist:"image-path"`
		Devices   []device `plist:"system-entities"`
	}

	type hdinfo struct {
		Images []image `plist:"images"`
	}

	pl, err := exec.Command("hdiutil", "info", "-plist").Output()
	if err != nil {
		return "", "", err
	}

	result := hdinfo{}
	_, err = plist.Unmarshal(pl, &result)
	if err != nil {
		return "", "", err
	}

	for _, image := range result.Images {
		if image.ImagePath == path {
			for _, device := range image.Devices {
				matched, err := regexp.MatchString("^/dev/disk[0-9]+$", device.DevEntry)
				if err != nil {
					return "", "", err
				}

				if matched {
					return device.DevEntry, strings.Replace(device.DevEntry, "disk", "rdisk", 1), nil
				}
			}
		}
	}

	return "", "", fmt.Errorf("Disk not attached")
}

func attachDisk(path string) error {
	_, _, err := findDisk(path)
	if err == nil || err.Error() != "Disk not attached" {
		return err
	}

	return exec.Command("hdiutil", "attach", "-nomount", "-noverify", "-noautofsck", path).Run()
}

func detachDisk(path string) error {
	dev, _, err := findDisk(path)
	if err != nil {
		return err
	}

	return exec.Command("hdiutil", "detach", dev).Run()
}
