package os

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/blang/semver"
)

type Version struct {
	Tag string `json:"tag_name"`
	Url string `json:"tarball_url"`
}

func getReleases() ([]semver.Version, map[string]Version, error) {
	res, err := http.Get("https://api.github.com/repos/nlf/dhyve-os/releases")
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()
	vers := []Version{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(vers)
	if err != nil {
		return nil, nil, err
	}

	versions := map[string]Version{}
	rng, err := semver.ParseRange(">= 3.0.0-beta0")
	if err != nil {
		return nil, nil, err
	}

	svers := []semver.Version{}
	for _, vers := range versions {
		v := semver.MustParse(vers.Tag)
		if rng(v) {
			versions[vers.Tag] = vers
			svers = append(svers, v)
		}
	}

	semver.Sort(svers)
	return svers, versions, nil
}

func Latest() (Version, error) {
	svers, mvers, err := getReleases()
	if err != nil {
		return Version{}, err
	}

	latest := svers[len(svers)-1]
	return mvers[latest.String()], nil
}

func Specific(vers string) (Version, error) {
	sv := semver.MustParse(vers)
	svers, mvers, err := getReleases()
	if err != nil {
		return Version{}, err
	}

	for _, v := range svers {
		if v.Equals(sv) {
			return mvers[v.String()], nil
		}
	}

	return Version{}, fmt.Errorf("Cannot find version %s", vers)
}

func Download(dir string, v Version) error {
	path := fmt.Sprintf("%d/%s.tar.gz", dir, v.Tag)
	output, err := os.Create(path)
	if err != nil {
		return err
	}

	res, err := http.Get(v.Url)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	length, err := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}

	_, err = io.CopyN(output, res.Body, length)
	if err != nil {
		return err
	}

	return exec.Command("tar", "xf", path).Run()
}
