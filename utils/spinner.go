package utils

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type Step struct {
	Prefix string
	Action func() error
}

type Steps []Step

func Spin(fm Steps) error {
	for _, fn := range fm {
		s := spinner.New(spinner.CharSets[9], time.Millisecond*100)
		s.Prefix = fmt.Sprintf("%s: ", fn.Prefix)
		s.Start()
		err := fn.Action()
		s.Stop()
		if err != nil {
			fmt.Printf("\r%s: ERROR!\n", fn.Prefix)
			return err
		} else {
			fmt.Printf("\r%s: done\n", fn.Prefix)
		}
	}

	return nil
}
