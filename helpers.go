package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

type User struct {
	Name string
	Home string
	Uid  int
	Gid  int
}

func getUser() User {
	currentUser, err := user.Current()
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	gid, err := strconv.Atoi(currentUser.Gid)
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	return User{
		Name: currentUser.Name,
		Home: currentUser.HomeDir,
		Uid:  uid,
		Gid:  gid,
	}
}

func getPath() string {
	currentUser := getUser()
	return filepath.Join(currentUser.Home, "._dlite")
}

func promptString(question, def string) (string, error) {
	prompt := fmt.Sprintf("%s:", question)
	if def != "" {
		prompt += fmt.Sprintf(" [%s]", def)
	}
	res, err := ui.Ask(prompt)
	if err != nil {
		return "", err
	}

	if res == "" {
		return def, nil
	}

	return res, nil
}

func promptInt(question string, def int) (int, error) {
	prompt := fmt.Sprintf("%s: [%d]", question, def)
	res, err := ui.Ask(prompt)
	if err != nil {
		return -1, err
	}

	if res == "" {
		return def, nil
	}

	return strconv.Atoi(res)
}
