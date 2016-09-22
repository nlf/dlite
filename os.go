package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/blang/semver"
)

var allowedRange = semver.MustParseRange(">=3.0.0-beta0")

type asset struct {
	Name string `json:"name"`
	Url  string `json:"browser_download_url"`
}

type version struct {
	Tag    string  `json:"tag_name"`
	Url    string  `json:"tarball_url"`
	Assets []asset `json:"assets"`
}

func getLatestOS() ([]asset, error) {
	res, err := http.Get("https://api.github.com/repos/nlf/dhyve-os/releases")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	versions := []version{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&versions)
	if err != nil {
		return nil, err
	}

	versionsMap := map[string]version{}
	availableVersions := []semver.Version{}
	for _, version := range versions {
		ver := semver.MustParse(strings.TrimPrefix(version.Tag, "v"))
		if allowedRange(ver) {
			versionsMap[ver.String()] = version
			availableVersions = append(availableVersions, ver)
		}
	}

	semver.Sort(availableVersions)
	return versionsMap[availableVersions[len(availableVersions)-1].String()].Assets, nil
}

func DownloadOS(target string) error {
	latest, err := getLatestOS()
	if err != nil {
		return err
	}

	for _, asset := range latest {
		tempPath := path.Join(target, asset.Name)
		output, err := os.Create(tempPath)
		if err != nil {
			return err
		}

		res, err := http.Get(asset.Url)
		defer res.Body.Close()
		if err != nil {
			return err
		}

		length, err := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			return err
		}

		_, err = io.CopyN(output, res.Body, length)
		if err != nil {
			return err
		}
	}

	return nil
}
