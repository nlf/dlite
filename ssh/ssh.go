package ssh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nlf/dlite/config"
)

var startMarker []byte = []byte("# begin dlite")
var endMarker []byte = []byte("# end dlite\n")

type SSH struct {
	config *config.Config
}

func (s *SSH) Generate() error {
	path := filepath.Join(s.config.Dir, "key")
	os.RemoveAll(path)
	os.RemoveAll(path + ".pub")

	return exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-C", "dlite", "-f", path, "-N", "").Run()
}

func (s *SSH) Configure(ip string) error {
	entry := fmt.Sprintf("# begin dlite\nHost %s\n  HostName %s\n  IdentityFile %s\n  User docker\n  StrictHostKeyChecking no\n# end dlite\n", s.config.Hostname, ip, filepath.Join(s.config.Dir, "key"))
	path := filepath.Join(s.config.Home, ".ssh", "config")
	os.MkdirAll(filepath.Dir(path), 0755)

	config, err := ioutil.ReadFile(filepath.Join(s.config.Home, ".ssh", "config"))
	if err != nil {
		if os.IsNotExist(err) {
			return ioutil.WriteFile(path, []byte(entry), 0644)
		}
		return err
	}

	begin := bytes.Index(config, startMarker)
	end := bytes.Index(config, endMarker)

	var temp []byte

	if begin > -1 && end > -1 {
		temp = append(config[:begin], config[end+len(endMarker):]...)
		temp = append(bytes.TrimSpace(temp), '\n')
	} else {
		temp = config
	}

	if len(temp) > 0 && !bytes.HasSuffix(temp, []byte("\n")) {
		temp = append(temp, []byte("\n")...)
	}
	temp = append(temp, []byte(entry)...)
	return ioutil.WriteFile(path, temp, 0644)
}

func (s *SSH) RemoveConfig() error {
	path := filepath.Join(s.config.Home, ".ssh", "config")
	config, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	begin := bytes.Index(config, startMarker)
	end := bytes.Index(config, endMarker)

	if begin == -1 && end == -1 {
		return nil
	}

	temp := append(config[:begin], config[end+len(endMarker):]...)
	temp = append(bytes.TrimSpace(temp), '\n')
	if len(temp) > 0 && !bytes.HasSuffix(temp, []byte("\n")) {
		temp = append(temp, []byte("\n")...)
	}
	return ioutil.WriteFile(path, temp, 0644)
}

func New(config *config.Config) *SSH {
	return &SSH{
		config: config,
	}
}
