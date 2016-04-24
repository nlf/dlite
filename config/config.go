package config

import (
	"bytes"
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

type Config struct {
	path          string `json:"-"`
	Home          string `json:"-"`
	Dir           string `json:"-"`
	Uid           int    `json:"-"`
	Gid           int    `json:"-"`
	Username      string `json:"-"`
	Uuid          string `json:"uuid"`
	CpuCount      int    `json:"cpu_count"`
	DiskSize      int    `json:"disk_size"`
	Memory        int    `json:"memory"`
	Hostname      string `json:"hostname"`
	DNSServer     string `json:"dns_server"`
	Extra         string `json:"extra"`
	DockerVersion string `json:"docker_version"`
	Route         bool   `json:"route"`
}

func (c *Config) Save() error {
	err := os.MkdirAll(c.Dir, 0755)
	if err != nil {
		return err
	}

	output, err := os.Create(c.path)
	if err != nil {
		return err
	}

	defer output.Close()

	buf, err := json.Marshal(c)
	if err != nil {
		return err
	}

	prettified := bytes.Buffer{}
	json.Indent(&prettified, buf, "", "  ")
	_, err = prettified.WriteTo(output)
	if err != nil {
		return err
	}

	return os.Chown(c.path, c.Uid, c.Gid)
}

func (c *Config) Load() error {
	file, err := os.Open(c.path)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(c)
}

func New(username string) (*Config, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return nil, err
	}

	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return nil, err
	}

	return &Config{
		path:     filepath.Join(u.HomeDir, ".dlite", "config.json"),
		Dir:      filepath.Join(u.HomeDir, ".dlite"),
		Home:     u.HomeDir,
		Uid:      uid,
		Gid:      gid,
		Username: username,
	}, nil
}
