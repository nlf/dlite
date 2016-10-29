package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/kardianos/osext"
	"github.com/urfave/cli"
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
		fmt.Println(err)
		os.Exit(1)
	}

	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gid, err := strconv.Atoi(currentUser.Gid)
	if err != nil {
		fmt.Println(err)
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
	return filepath.Join(user.Home, ".dlite")
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

func statusRequest() (*VMStatus, *cli.ExitError) {
	user := getUser()
	req, err := http.NewRequest("GET", "http://127.0.0.1:1050/status", nil)
	if err != nil {
		return nil, cli.NewExitError(err.Error(), 1)
	}

	req.Header.Add("X-Username", user.Name)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, cli.NewExitError(getRequestError(err), 1)
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		status := VMStatusError{}
		err := decoder.Decode(&status)
		if err != nil {
			return nil, cli.NewExitError(err.Error(), 1)
		}

		return nil, cli.NewExitError(status.Message, 1)
	}

	status := VMStatus{}
	err = decoder.Decode(&status)
	if err != nil {
		return nil, cli.NewExitError(err.Error(), 1)
	}

	return &status, nil
}

func stringRequest(action string) *cli.ExitError {
	user := getUser()
	req, err := http.NewRequest("POST", fmt.Sprintf("http://127.0.0.1:1050/%s", action), nil)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	req.Header.Add("X-Username", user.Name)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	code := 0
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		code = 1
	}

	return cli.NewExitError(string(body), code)
}

func runSetup(hostname, home string) *cli.ExitError {
	exe, err := osext.Executable()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	output, err := exec.Command("sudo", exe, "setup", "--hostname", hostname, "--home", home).Output()
	code := 0
	if err != nil {
		code = 1
	}
	return cli.NewExitError(string(output), code)
}

func runCleanup(hostname, home string) *cli.ExitError {
	exe, err := osext.Executable()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	output, err := exec.Command("sudo", exe, "cleanup", "--hostname", hostname, "--home", home).Output()
	code := 0
	if err != nil {
		code = 1
	}
	return cli.NewExitError(string(output), code)
}

func ensureRoot() *cli.ExitError {
	if uid := os.Geteuid(); uid != 0 {
		return cli.NewExitError("This command requires sudo", 1)
	}

	if uid := os.Getenv("SUDO_UID"); uid == "" {
		return cli.NewExitError("This command requires sudo", 1)
	}

	return nil
}
