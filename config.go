package main

import (
	"path/filepath"

	"github.com/tsuru/config"
)

type Config struct {
	Id       string
	Hostname string
	Disk     int
	DiskPath string
	Cpu      int
	Memory   int
	DNS      string
	Docker   string
	Extra    string
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

	return config.WriteConfigFile(configFile, 0644)
}
