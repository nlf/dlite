package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func getUser() *user.User {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Unable to lookup current user")
		os.Exit(1)
	}

	return currentUser
}

func getPath() string {
	currentUser := getUser()
	return filepath.Join(currentUser.HomeDir, "._dlite")
}
