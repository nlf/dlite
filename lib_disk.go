package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/DHowett/go-plist"
)

type Device struct {
	DevEntry string `plist:"dev-entry"`
}

type Image struct {
	ImagePath string   `plist:"image-path"`
	Devices   []Device `plist:"system-entities"`
}

type HDIInfo struct {
	Images []Image `plist:"images"`
}

func changePermissions(path string) error {
	var uid, gid int
	var err error

	suid := os.Getenv("SUDO_UID")
	if suid != "" {
		uid, err = strconv.Atoi(suid)
		if err != nil {
			return err
		}
	} else {
		uid = os.Getuid()
	}

	sgid := os.Getenv("SUDO_GID")
	if sgid != "" {
		gid, err = strconv.Atoi(sgid)
		if err != nil {
			return err
		}
	} else {
		gid = os.Getgid()
	}

	return os.Chown(path, uid, gid)
}

func CreateDir() error {
	path := os.ExpandEnv("$HOME/.dlite")

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	return changePermissions(path)
}

func RemoveDir() error {
	path := os.ExpandEnv("$HOME/.dlite")
	return os.RemoveAll(path)
}

func GetDiskPath() (string, error) {
	pl, err := exec.Command("hdiutil", "info", "-plist").Output()
	if err != nil {
		return "", err
	}

	var results HDIInfo
	_, err = plist.Unmarshal(pl, &results)
	if err != nil {
		return "", err
	}

	for _, image := range results.Images {
		if image.ImagePath == os.ExpandEnv("$HOME/.dlite/disk.sparseimage") {
			for _, device := range image.Devices {
				matched, err := regexp.MatchString("^/dev/disk[0-9]+$", device.DevEntry)
				if err != nil {
					return "", err
				}

				if matched {
					return strings.Replace(device.DevEntry, "disk", "rdisk", 1), nil
				}
			}
		}
	}

	return "", fmt.Errorf("Disk not attached")
}

func AttachDisk() (string, error) {
	_, err := GetDiskPath()
	if err != nil {
		if err.Error() != "Disk not attached" {
			return "", err
		}

		err = exec.Command("hdiutil", "attach", "-nomount", "-noverify", "-noautofsck", os.ExpandEnv("$HOME/.dlite/disk.sparseimage")).Run()
		if err != nil {
			return "", err
		}
	}

	return GetDiskPath()
}

func DetachDisk() error {
	path, err := GetDiskPath()
	if err != nil {
		return err
	}

	return exec.Command("hdiutil", "detach", strings.Replace(path, "rdisk", "disk", 1)).Run()
}

func CreateDisk(size int) error {
	DetachDisk()
	path := os.ExpandEnv("$HOME/.dlite/disk")
	err := os.RemoveAll(path + ".sparseimage")
	if err != nil {
		return err
	}

	err = exec.Command("hdiutil", "create", "-size", fmt.Sprintf("%dg", size), "-type", "SPARSE", "-layout", "MBRSPUD", path).Run()
	if err != nil {
		return err
	}

	return changePermissions(path + ".sparseimage")
}
