package nfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nlf/dlite/config"
)

type NFS struct {
	config *config.Config
}

func (n *NFS) AddExport() error {
	export := fmt.Sprintf("%s -network 192.168.64.0 -mask 255.255.255.0 -alldirs -maproot=root:wheel  # added by dlite", n.config.Home)

	exports, err := ioutil.ReadFile("/etc/exports")
	if err != nil {
		if os.IsNotExist(err) {
			return ioutil.WriteFile("/etc/exports", []byte(export), 0644)
		}
		return err
	}

	if strings.Contains(string(exports), "# added by dlite") {
		err := n.RemoveExport()
		if err != nil {
			return err
		}
	}

	exports = append(exports, []byte("\n"+export)...)
	return ioutil.WriteFile("/etc/exports", exports, 0644)
}

func (n *NFS) RemoveExport() error {
	exports, err := ioutil.ReadFile("/etc/exports")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	re := regexp.MustCompile("(?m)[\r\n]+^.*# added by dlite.*$")
	removed := re.ReplaceAllString(string(exports), "")
	return ioutil.WriteFile("/etc/exports", []byte(removed), 0644)
}

func (n *NFS) Start() error {
	output, err := exec.Command("nfsd", "status").Output()
	if err != nil {
		return err
	}

	if strings.Contains(string(output), "nfsd is running") {
		return n.Reload()
	}

	return exec.Command("nfsd", "start").Run()
}

func (n *NFS) Reload() error {
	output, err := exec.Command("nfsd", "status").Output()
	if err != nil {
		return err
	}

	if !strings.Contains(string(output), "nfsd is running") {
		return n.Start()
	}

	return exec.Command("nfsd", "update").Run()
}

func New(config *config.Config) *NFS {
	return &NFS{
		config: config,
	}
}
