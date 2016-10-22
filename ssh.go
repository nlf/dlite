package main

import (
	"fmt"
	"os/exec"
)

func generateKeys(user User) error {
	base := getPath(user)
	return exec.Command("ssh-keygen", "-f", fmt.Sprintf("%s/key", base), "-P", "\"\"").Run()
}
