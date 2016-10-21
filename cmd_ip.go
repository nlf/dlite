package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/cli"
)

type ipCommand struct{}

func (c *ipCommand) Run(args []string) int {
	user := getUser()

	req, err := http.NewRequest("GET", "http://127.0.0.1:1050/status", nil)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	req.Header.Add("X-Username", user.Name)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		status := VMStatusError{}
		json.Unmarshal(body, &status)
		ui.Error(status.Message)
		return 1
	}

	status := VMStatus{}
	err = json.Unmarshal(body, &status)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	ui.Output(status.IP)
	return 0
}

func (c *ipCommand) Synopsis() string {
	return "get the ip of the virtual machine"
}

func (c *ipCommand) Help() string {
	return "returns the ip address of the virtual machine"
}

func ipFactory() (cli.Command, error) {
	return &ipCommand{}, nil
}
