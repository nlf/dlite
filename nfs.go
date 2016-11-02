package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func ensureNFS(home string) error {
	addr, err := getNetAddress()
	if err != nil {
		return err
	}

	mask, _ := getNetMask()
	export := fmt.Sprintf("%s -network %s -mask %s -alldirs -maproot=root:wheel", home, addr, mask)

	if _, err = os.Stat("/etc/exports"); os.IsNotExist(err) {
		err := ioutil.WriteFile("/etc/exports", []byte(""), 0644)
		if err != nil {
			return err
		}
	}

	rawExports, err := ioutil.ReadFile("/etc/exports")
	if err != nil {
		return err
	}

	needsExport := true
	for _, line := range strings.Split(string(rawExports), "\n") {
		if strings.HasPrefix(line, export) {
			needsExport = false
			break
		}
	}

	if needsExport {
		file, err := os.OpenFile("/etc/exports", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		defer file.Close()

		_, err = file.WriteString(export + "\n")
		if err != nil {
			return err
		}
	}

	output, err := exec.Command("nfsd", "checkexports").Output()
	if err != nil {
		return fmt.Errorf("There was a problem updating the /etc/exports file, please resolve the issue and run 'sudo nfsd restart'\n%s", string(output))
	}

	output, _ = exec.Command("nfsd", "status").Output()
	enabled := false
	running := false
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "is enabled") {
			enabled = true
		} else if strings.Contains(line, "is running") {
			running = true
		}
	}

	if !enabled {
		output, err = exec.Command("nfsd", "enable").Output()
	} else if !running {
		output, err = exec.Command("nfsd", "start").Output()
	} else {
		output, err = exec.Command("nfsd", "restart").Output()
	}

	if err != nil {
		return fmt.Errorf(string(output))
	}
	return nil
}

func removeNFS(home string) error {
	addr, err := getNetAddress()
	if err != nil {
		return err
	}

	mask, _ := getNetMask()
	export := fmt.Sprintf("%s -network %s -mask %s -alldirs -maproot=root:wheel", home, addr, mask)

	rawExports, err := ioutil.ReadFile("/etc/exports")
	if err != nil {
		return err
	}

	exportMatcher := regexp.MustCompile(fmt.Sprintf("(?m)^%s\n?$", export))
	newExports := exportMatcher.ReplaceAllString(string(rawExports), "")

	err = ioutil.WriteFile("/etc/exports", []byte(newExports), 0644)
	if err != nil {
		return err
	}

	output, err := exec.Command("nfsd", "restart").Output()
	if err != nil {
		return fmt.Errorf(string(output))
	}

	return nil
}
