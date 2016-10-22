package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/urfave/cli"
)

func ask(question string) string {
	var answer string
	fmt.Printf("%s: ", question)
	fmt.Scanln(&answer)
	return answer
}

func askString(question, def string) string {
	prompt := fmt.Sprintf("%s:", question)
	if def != "" {
		prompt += fmt.Sprintf(" [%s]", def)
	}
	res := ask(prompt)
	if res == "" {
		return def
	}

	return res
}

func askInt(question string, def int) int {
	prompt := fmt.Sprintf("%s: [%d]", question, def)
	res := ask(prompt)
	if res == "" {
		return def
	}

	i, err := strconv.Atoi(res)
	if err != nil {
		return askInt(question, def)
	}

	return i
}

func confirm(question string) bool {
	answer := ask(question)

	switch strings.ToLower(answer) {
	case "yes", "y":
		return true
	case "no", "n":
		return false
	default:
		return confirm(question)
	}
}

func spin(prefix string, f func() error) *cli.ExitError {
	spin := spinner.New(spinner.CharSets[9], time.Millisecond*100)
	spin.Prefix = fmt.Sprintf("%s: ", prefix)
	spin.Start()
	err := f()
	spin.Stop()
	if err != nil {
		cliError, ok := err.(*cli.ExitError)
		if !ok || cliError.ExitCode() > 0 {
			fmt.Printf("\r%s: ERROR!\n", prefix)
		} else {
			fmt.Printf("\r%s: done\n", prefix)
		}
	} else {
		fmt.Printf("\r%s: done\n", prefix)
	}

	if err != nil {
		cliError, ok := err.(*cli.ExitError)
		if ok {
			return cliError
		}

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
