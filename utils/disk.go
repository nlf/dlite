package utils

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
)

func changePermissions(path string) error {
	var uid, gid int
	var err error

	suid := os.Getenv("SUDO_UID")
	if suid != "" {
		uid, err = strconv.Atoi(suid)
		if err != nil {
			return err
		}
	} else {
		uid = os.Getuid()
	}

	sgid := os.Getenv("SUDO_GID")
	if sgid != "" {
		gid, err = strconv.Atoi(sgid)
		if err != nil {
			return err
		}
	} else {
		gid = os.Getgid()
	}

	return os.Chown(path, uid, gid)
}

func CreateDir() error {
	path := os.ExpandEnv("$HOME/.dlite")

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	return changePermissions(path)
}

func RemoveDir() error {
	path := os.ExpandEnv("$HOME/.dlite")
	return os.RemoveAll(path)
}

func CreateDisk(sshKey string, size int) error {
	if strings.Contains(sshKey, "$HOME") {
		username := os.Getenv("SUDO_USER")
		if username == "" {
			username = os.Getenv("USER")
		}

		me, err := user.Lookup(username)
		if err != nil {
			return err
		}

		sshKey = strings.Replace(sshKey, "$HOME", me.HomeDir, -1)
	}

	sshKey = os.ExpandEnv(sshKey)
	keyBytes, err := ioutil.ReadFile(sshKey)
	if err != nil {
		return err
	}

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
			return err
		}

		if _, err = tarball.Write(file.Body); err != nil {
			return err
		}
	}

	if err = tarball.Close(); err != nil {
		return err
	}

	path := os.ExpandEnv("$HOME/.dlite/disk.img")
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	_, err = f.Seek(int64(size*1073741824-1), 0)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte{0})
	if err != nil {
		return err
	}

	return changePermissions(path)
}
