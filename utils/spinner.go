package utils

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type wrapped func() error
type FunctionMap map[string]wrapped

func Spin(fm FunctionMap) error {
	for prefix, fn := range fm {
		s := spinner.New(spinner.CharSets[9], time.Millisecond*100)
		s.Prefix = fmt.Sprintf("%s: ", prefix)
		s.Start()
		err := fn()
		s.Stop()
		if err != nil {
			fmt.Printf("\r%s: ERROR - %s\n", prefix, err)
			return err
		} else {
			fmt.Printf("\r%s: done\n", prefix)
		}
	}

	return nil
}
