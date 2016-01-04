package main

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/johanneswuerbach/nfsexports"
	"github.com/satori/go.uuid"
)

func installHandler(ctx *cli.Context) {
	if uid := os.Geteuid(); uid != 0 {
		log.Fatalln("the install command requires 'sudo'")
	}

	// read params
	cpuCount := ctx.Int("cpu")
	diskSize := ctx.Int("disk")
	memSize := ctx.Int("memory")
	sshKey := os.ExpandEnv(ctx.String("ssh-key"))

	log.Printf("Installing with %dGB disk, %dGB of ram and %d CPUs\nUsing SSH key from %s\n", diskSize, memSize, cpuCount, sshKey)
	uuid := uuid.NewV4().String()

	createDisk(sshKey, diskSize)
	downloadDhyveOS()
	addExport(uuid)
	saveConfig(uuid, cpuCount, memSize)
}

func createDisk(sshKey string, diskSize int) {
	// read ssh key file
	keyBytes, err := ioutil.ReadFile(sshKey)
	if err != nil {
		log.Fatalf("Failed to read SSH key file at %s\n%s", sshKey, err.Error())
	}

	// create the tarball header
	buffer := new(bytes.Buffer)
	tarball := tar.NewWriter(buffer)
	files := []struct {
		Name string
		Body []byte
	}{
		{"dhyve, please format-me", []byte("dhyve, please format-me")},
		{".ssh/authorized_keys", keyBytes},
	}

	for _, file := range files {
		if err = tarball.WriteHeader(&tar.Header{
			Name: file.Name,
			Mode: 0644,
			Size: int64(len(file.Body)),
		}); err != nil {
			log.Fatalln(err)
		}

		if _, err = tarball.Write(file.Body); err != nil {
			log.Fatalln(err)
		}
	}

	if err = tarball.Close(); err != nil {
		log.Fatalln(err)
	}

	// write the tarball to a real file
	f, err := os.Create(filepath.Join(PREFIX, "/var/db/dlite/disk.img"))
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()
	_, err = f.Write(buffer.Bytes())
	if err != nil {
		log.Fatalln(err)
	}

	// write zeroes to the file
	halfGig := make([]byte, 536870912, 536870912)
	for count := 0; count < diskSize*2; count++ {
		_, err = f.Write(halfGig)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func downloadDhyveOS() {
	resp, err := http.Get("https://api.github.com/repos/nlf/dhyve-os/releases/latest")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	var latest struct {
		Tag string `json:"tag_name"`
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&latest)
	if err != nil {
		log.Fatalln(err)
	}

	files := []string{"bzImage", "rootfs.cpio.xz"}
	for _, file := range files {
		output, err := os.Create(filepath.Join(PREFIX, "/usr/share/dlite", file))
		if err != nil {
			log.Fatalln(err)
		}

		defer output.Close()

		resp, err = http.Get("https://github.com/nlf/dhyve-os/releases/download/" + latest.Tag + "/" + file)
		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()
		_, err = io.Copy(output, resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func addExport(uuid string) {
	export := fmt.Sprintf("/Users %s -alldirs -mapall=%s:%s", "-network 192.168.64.0 -mask 255.255.255.0", os.Getuid(), os.Getgid())
	if _, err := nfsexports.Add("", fmt.Sprintf("dlite %s", uuid), export); err != nil {
		log.Fatalln(err)
	}

	err := nfsexports.ReloadDaemon()
	if err != nil {
		log.Fatalln(err)
	}
}
