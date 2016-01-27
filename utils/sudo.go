package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/osext"
)

func AddSudoer() error {
	user := os.Getenv("SUDO_USER")

	path, err := osext.Executable()
	if err != nil {
		return err
	}
	sudoer := fmt.Sprintf("%s ALL=(ALL) NOPASSWD: /sbin/nfsd,%s\n", user, path)

	file, err := os.OpenFile("/private/etc/sudoers", os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	sudoers, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	lines := strings.Split(string(sudoers), "\n")
	exists := false
	for _, line := range lines {
		if line == sudoer {
			exists = true
			break
		}
	}

	if !exists {
		lines = append(lines, sudoer)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(strings.Join(lines, "\n")))
	return err
}

func RemoveSudoer() error {
	user := os.Getenv("SUDO_USER")

	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}
	sudoer := fmt.Sprintf("%s ALL=(ALL) NOPASSWD: /sbin/nfsd,%s", user, path)

	file, err := os.OpenFile("/private/etc/sudoers", os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	sudoers, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	lines := strings.Split(string(sudoers), "\n")
	for i, line := range lines {
		if line == sudoer {
			lines = append(lines[:i], lines[i+1:]...)
			break
		}
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	n, err := file.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return err
	}

	return file.Truncate(int64(n))
}

func EnsureSudo() {
	if uid := os.Geteuid(); uid != 0 {
		fmt.Println("This command requires sudo")
		os.Exit(1)
	}

	if uid := os.Getenv("SUDO_UID"); uid == "" {
		fmt.Println("This command requires sudo, please do not run it as root")
		os.Exit(1)
	}
}
