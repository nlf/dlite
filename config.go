package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Uuid     string `json:"uuid"`
	CpuCount int    `json:"cpu_count"`
	Memory   int    `json:"memory"`
}

func saveConfig(uuid string, cpus, mem int) {
	output, err := os.Create(filepath.Join(PREFIX, "/etc/dlite.conf"))
	if err != nil {
		log.Fatalln(err)
	}

	defer output.Close()
	config := Config{
		Uuid:     uuid,
		CpuCount: cpus,
		Memory:   mem,
	}

	b, err := json.Marshal(config)
	if err != nil {
		log.Fatalln(err)
	}

	output.Write(b)
}

func readConfig() Config {
	var config Config
	file, err := os.Open(filepath.Join(PREFIX, "/etc/dlite.conf"))
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}

	return config
}
