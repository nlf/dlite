package disk

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DHowett/go-plist"
	"github.com/nlf/dlite/config"
)

type Disk struct {
	config *config.Config
	path   string
}

type Device struct {
	DevEntry string `plist:"dev-entry"`
}

type Image struct {
	ImagePath string   `plist:"image-path"`
	Devices   []Device `plist:"system-entities"`
}

type HDInfo struct {
	Images []Image `plist:"images"`
}

func (d *Disk) Create() error {
	err := os.MkdirAll(filepath.Dir(d.path), 0755)
	if err != nil {
		return err
	}

	err = os.RemoveAll(d.path)
	if err != nil {
		return err
	}

	err = exec.Command("hdiutil", "create", "-size", fmt.Sprintf("%dg", d.config.DiskSize), "-type", "SPARSE", "-layout", "MBRSPUD", d.path).Run()
	if err != nil {
		return err
	}

	return os.Chown(d.path, d.config.Uid, d.config.Gid)
}

func (d *Disk) Attach() error {
	_, err := d.device()
	if err == nil || err.Error() != "Disk not attached" {
		return err
	}

	return exec.Command("hdiutil", "attach", "-nomount", "-noverify", "-noautofsck", d.path).Run()
}

func (d *Disk) Detach() error {
	dev, err := d.device()
	if err != nil {
		return err
	}

	return exec.Command("hdiutil", "detach", dev).Run()
}

func (d *Disk) device() (string, error) {
	pl, err := exec.Command("hdiutil", "info", "-plist").Output()
	if err != nil {
		return "", err
	}

	results := HDInfo{}
	_, err = plist.Unmarshal(pl, &results)
	if err != nil {
	}

	for _, image := range results.Images {
		if image.ImagePath == d.path {
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

func (d *Disk) Device() (string, error) {
	return d.device()
}

func (d *Disk) RawDevice() (string, error) {
	dev, err := d.device()
	if err != nil {
		return "", err
	}

	return strings.Replace(dev, "disk", "rdisk", 1), nil
}

func New(config *config.Config) *Disk {
	return &Disk{
		config: config,
		path:   filepath.Join(config.Dir, "disk.sparseimage"),
	}
}
