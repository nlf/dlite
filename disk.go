package main

import (
	"fmt"
	"os/exec"
)

func buildDisk(bin, path string, size, uid, gid int) error {
	return exec.Command(bin, "create", fmt.Sprintf("--size=%dGiB", size), path).Run()
}
