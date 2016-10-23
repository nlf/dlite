package main

import (
	"fmt"
	"os/exec"
)

func generateKeys(user User) error {
	base := getPath(user)
	output, err := exec.Command("ssh-keygen", "-f", fmt.Sprintf("%s/key", base), "-P", "").CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(output))
	}

	return nil
}
