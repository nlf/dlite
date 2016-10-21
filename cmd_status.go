package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/cli"
)

type statusCommand struct{}

func (c *statusCommand) Run(args []string) int {
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
		ui.Error(getRequestError(err))
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

	if status.Started {
		ui.Output("vm_state:       started")
		ui.Output(fmt.Sprintf("ip_address:     %s", status.IP))
		ui.Output(fmt.Sprintf("pid:            %d", status.Pid))
	} else {
		ui.Output("vm_state:       stopped")
	}
	ui.Output(fmt.Sprintf("id:             %s", status.Id))
	ui.Output(fmt.Sprintf("hostname:       %s", status.Hostname))
	ui.Output(fmt.Sprintf("disk_size:      %d", status.Disk))
	ui.Output(fmt.Sprintf("disk_path:      %s", status.DiskPath))
	ui.Output(fmt.Sprintf("cpu_cores:      %d", status.Cpu))
	ui.Output(fmt.Sprintf("memory:         %d", status.Memory))
	ui.Output(fmt.Sprintf("dns_server:     %s", status.DNS))
	ui.Output(fmt.Sprintf("docker_version: %s", status.Docker))
	ui.Output(fmt.Sprintf("docker_args:    %s", status.Extra))

	return 0
}

func (c *statusCommand) Synopsis() string {
	return "get the status of the virtual machine"
}

func (c *statusCommand) Help() string {
	return "returns the full status of the virtual machine, including configuration parameters"
}

func statusFactory() (cli.Command, error) {
	return &statusCommand{}, nil
}
