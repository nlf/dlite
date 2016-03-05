package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetLatestOSVersion() (string, error) {
	return "v2.2.0", nil
}

func DownloadOS(version string) error {
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	files := []string{"bzImage", "rootfs.cpio.xz"}
	for _, file := range files {
		path := os.ExpandEnv("$HOME/.dlite/" + file)
		output, err := os.Create(path)
		if err != nil {
			return err
		}

		defer output.Close()

		resp, err := http.Get("https://github.com/nlf/dhyve-os/releases/download/" + version + "/" + file)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		io.Copy(output, resp.Body)
		err = changePermissions(path)
		if err != nil {
			return err
		}
	}

	return nil
}
