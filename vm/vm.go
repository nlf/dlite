package vm

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/nlf/dlite/config"
)

type result struct {
	value string
	err   error
}

type VM struct {
	config *config.Config
	tty    string
	err    error
}

func (v *VM) IP() (string, error) {
	value := make(chan result, 1)
	matchRe := regexp.MustCompile(`.*name=([a-fA-F0-9\-]+)\n.*ip_address=([0-9\.]+)`)

	go func(matchRe *regexp.Regexp, config *config.Config) {
		attempts := 0

		for {
			if attempts >= 15 {
				value <- result{"", fmt.Errorf("Timed out waiting for IP address")}
				break
			}

			time.Sleep(time.Second)

			file, err := os.Open("/var/db/dhcpd_leases")
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}

				value <- result{"", err}
				break
			}

			defer file.Close()
			leases, err := ioutil.ReadAll(file)
			if err != nil {
				value <- result{"", err}
				break
			}

			matches := matchRe.FindAllStringSubmatch(string(leases), -1)
			for _, match := range matches {
				if match[1] == config.Uuid {
					value <- result{match[2], nil}
					break
				}
			}

			attempts++
		}
	}(matchRe, v.config)

	res := <-value
	return res.value, res.err
}

func New(config *config.Config) *VM {
	return &VM{
		config: config,
	}
}
