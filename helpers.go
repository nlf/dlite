package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
)

type User struct {
	Name string
	Home string
	Uid  int
	Gid  int
}

func lookupUser(username string) (*User, error) {
	rawUser, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}

	uid, err := strconv.Atoi(rawUser.Uid)
	if err != nil {
		return nil, err
	}

	gid, err := strconv.Atoi(rawUser.Gid)
	if err != nil {
		return nil, err
	}

	return &User{
		Name: rawUser.Username,
		Home: rawUser.HomeDir,
		Uid:  uid,
		Gid:  gid,
	}, nil
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
		Name: currentUser.Username,
		Home: currentUser.HomeDir,
		Uid:  uid,
		Gid:  gid,
	}
}

func getPath(user User) string {
	return filepath.Join(user.Home, "._dlite")
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

func spin(prefix string, f func() error) error {
	spin := spinner.New(spinner.CharSets[9], time.Millisecond*100)
	spin.Prefix = fmt.Sprintf("%s: ", prefix)
	spin.Start()
	err := f()
	spin.Stop()
	if err != nil {
		fmt.Printf("\r%s: ERROR!\n", prefix)
	} else {
		fmt.Printf("\r%s: done\n", prefix)
	}
	return err
}

func getRequestError(err error) string {
	urlError, ok := err.(*url.Error)
	if ok {
		netError, ok := urlError.Err.(*net.OpError)
		if ok {
			sysError, ok := netError.Err.(*os.SyscallError)
			if ok {
				if sysError.Err == syscall.ECONNREFUSED {
					return "Connection refused - is the dlite daemon running?"
				}
			}

			if netError.Timeout() {
				return "Request timed out, please try again"
			}
		}
	}

	return err.Error()
}
