package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetLatestOSVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/nlf/dhyve-os/releases/latest")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	var latest struct {
		Tag string `json:"tag_name"`
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&latest)
	if err != nil {
		return "", err
	}

	return latest.Tag, nil
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

		url := "https://github.com/nlf/dhyve-os/releases/download/" + version + "/" + file
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		contentLength, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			return fmt.Errorf("Unable to get Content-Length of '%s': %s", url, err.Error())
		}

		if _, err := io.CopyN(output, resp.Body, contentLength); err != nil {
			return fmt.Errorf("Unable to download '%s': %s", url, err.Error())
		}

		if err := changePermissions(path); err != nil {
			return err
		}
	}

	return nil
}
