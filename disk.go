package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"

	"github.com/DHowett/go-plist"
)

type Disk struct {
	Path string
	Size int
	uid  int
	gid  int
}

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

func NewDisk(path string, size int, owner *user.User) (*Disk, error) {
	gid, err := strconv.Atoi(owner.Gid)
	if err != nil {
		return nil, err
	}

	uid, err := strconv.Atoi(owner.Uid)
	if err != nil {
		return nil, err
	}

	return &Disk{
		Path: path,
		Size: size,
		gid:  gid,
		uid:  uid,
	}, nil
}

func (d *Disk) Build() error {
	err := exec.Command("hdiutil", "create", "-size", fmt.Sprintf("%dg", d.Size), "-type", "SPARSE", "-layout", "MBRSPUD", d.Path).Run()
	if err != nil {
		return err
	}

	return os.Chown(d.Path, d.uid, d.gid)
}

func (d *Disk) Attach() error {
	_, err := d.Device()
	if err == nil || err.Error() != "Disk not attached" {
		return err
	}

	return exec.Command("hdiutil", "attach", "-nomount", "-noverify", "-noautofsck", d.Path).Run()
}

func (d *Disk) Detach() error {
	dev, err := d.Device()
	if err != nil {
		return err
	}

	return exec.Command("hdiutil", "detach", dev).Run()
}

func (d *Disk) Device() (string, error) {
	pl, err := exec.Command("hdiutil", "info", "-plist").Output()
	if err != nil {
		return "", err
	}

	result := hdinfo{}
	_, err = plist.Unmarshal(pl, &result)
	if err != nil {
		return "", err
	}

	for _, image := range result.Images {
		if image.ImagePath == d.Path {
			for _, device := range image.Devices {
				matched, err := regexp.MatchString("^/dev/disk[0-9]+$", device.DevEntry)
				if err != nil {
					return "", err
				}

				if matched {
					return device.DevEntry, nil
				}
			}
		}
	}

	return "", fmt.Errorf("Disk not attached")
}

func (d *Disk) RawDevice() (string, error) {
	dev, err := d.Device()
	if err != nil {
		return "", err
	}

	return strings.Replace(dev, "disk", "rdisk", 1), nil
}
