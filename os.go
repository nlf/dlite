package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/blang/semver"
)

var allowedRange = semver.MustParseRange(">=1.0.0-beta0")

type File struct {
	Name string `json:"name"`
	Url  string `json:"browser_download_url"`
}

type Release struct {
	Version semver.Version
	Tag     string `json:"tag_name"`
	Files   []File `json:"assets"`
}
type Releases []Release

func (rs Releases) Len() int {
	return len(rs)
}

func (rs Releases) Less(i, j int) bool {
	return rs[i].Version.LT(rs[j].Version)
}

func (rs Releases) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func getOSReleases() (Releases, error) {
	res, err := http.Get("https://api.github.com/repos/nlf/dlite-os/releases")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	releases := Releases{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&releases)
	if err != nil {
		return nil, err
	}

	allowedReleases := Releases{}
	for _, release := range releases {
		release.Version = semver.MustParse(strings.TrimPrefix(release.Tag, "v"))
		if allowedRange(release.Version) {
			allowedReleases = append(allowedReleases, release)
		}
	}

	sort.Sort(allowedReleases)
	return allowedReleases, nil
}

func getLatestOSRelease() (Release, error) {
	releases, err := getOSReleases()
	if err != nil {
		return Release{}, err
	}

	return releases[len(releases)-1], nil
}

func downloadOS(target string) error {
	latest, err := getLatestOSRelease()
	if err != nil {
		return err
	}

	for _, file := range latest.Files {
		tempPath := path.Join(target, file.Name)
		output, err := os.Create(tempPath)
		if err != nil {
			return err
		}

		res, err := http.Get(file.Url)
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
