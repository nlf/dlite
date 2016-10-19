package main

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/cli"
)

type stopCommand struct{}

func (c *stopCommand) Run(args []string) int {
	user := getUser()
	var response string

	err := spin("Stopping the virtual machine", func() error {
		req, err := http.NewRequest("POST", "http://127.0.0.1:1050/stop", nil)
		if err != nil {
			return err
		}

		req.Header.Add("X-Username", user.Name)

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		response = string(body)
		if res.StatusCode < 200 || res.StatusCode >= 400 {
			return errors.New(response)
		}

		return nil
	})

	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *stopCommand) Synopsis() string {
	return "stop the virtual machine"
}

func (c *stopCommand) Help() string {
	return "this command will signal the virtual machine to stop"
}

func stopFactory() (cli.Command, error) {
	return &stopCommand{}, nil
}
