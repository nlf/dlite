package docker

import (
	"encoding/json"
	"net/http"
	"strings"
)

func Latest() (string, error) {
	res, err := http.Get("https://api.github.com/repos/docker/docker/releases/latest")
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	var latest struct {
		Tag string `json:"tag_name"`
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&latest)
	if err != nil {
		return "", nil
	}

	return strings.TrimPrefix(latest.Tag, "v"), nil
}
