package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Uuid          string `json:"uuid"`
	CpuCount      int    `json:"cpu_count"`
	DiskSize      int    `json:"disk_size"`
	Memory        int    `json:"memory"`
	Hostname      string `json:"hostname"`
	Share         string `json:"share"`
	DNSServer     string `json:"dns_server"`
	Extra         string `json:"extra"`
	DockerVersion string `json:"docker_version"`
}

func SaveConfig(config Config) error {
	path := os.ExpandEnv("$HOME/.dlite/config.json")
	output, err := os.Create(path)
	if err != nil {
		return err
	}

	defer output.Close()

	b, err := json.Marshal(config)
	if err != nil {
		return err
	}

	output.Write(b)
	return changePermissions(path)
}

func ReadConfig() (Config, error) {
	var config Config
	file, err := os.Open(os.ExpandEnv("$HOME/.dlite/config.json"))
	if err != nil {
		return Config{}, err
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
