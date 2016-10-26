package main

import (
	"path/filepath"

	"github.com/tsuru/config"
)

type Config struct {
	Id       string `json:"id"`
	Hostname string `json:"hostname"`
	Disk     int    `json:"disk_size"`
	DiskPath string `json:"disk_path"`
	Cpu      int    `json:"cpu_cores"`
	Memory   int    `json:"memory"`
	DNS      string `json:"dns_server"`
	Docker   string `json:"docker_version"`
	Extra    string `json:"docker_args"`
	Route    bool   `json:"route"`
}

func readConfig(path string) (Config, error) {
	cfg := Config{}
	configFile := filepath.Join(path, "config.yaml")

	err := config.ReadConfigFile(configFile)
	if err != nil {
		return cfg, err
	}

	cfg.Id, err = config.GetString("id")
	if err != nil {
		return cfg, err
	}

	cfg.Hostname, err = config.GetString("hostname")
	if err != nil {
		return cfg, err
	}

	cfg.DiskPath = filepath.Join(path, "disk.qcow")
	cfg.Disk, err = config.GetInt("disk")
	if err != nil {
		return cfg, err
	}

	cfg.Cpu, err = config.GetInt("cpu")
	if err != nil {
		return cfg, err
	}

	cfg.Memory, err = config.GetInt("memory")
	if err != nil {
		return cfg, err
	}

	cfg.DNS, err = config.GetString("dns")
	if err != nil {
		return cfg, err
	}

	cfg.Docker, err = config.GetString("docker")
	if err != nil {
		return cfg, err
	}

	cfg.Extra, err = config.GetString("extra")
	if err != nil {
		return cfg, err
	}

	cfg.Route, err = config.GetBool("route")
	return cfg, err
}

func writeConfig(path string, cfg Config) error {
	configFile := filepath.Join(path, "config.yaml")

	config.Set("id", cfg.Id)
	config.Set("hostname", cfg.Hostname)
	config.Set("disk", cfg.Disk)
	config.Set("cpu", cfg.Cpu)
	config.Set("memory", cfg.Memory)
	config.Set("dns", cfg.DNS)
	config.Set("docker", cfg.Docker)
	config.Set("extra", cfg.Extra)
	config.Set("route", cfg.Route)

	return config.WriteConfigFile(configFile, 0644)
}
