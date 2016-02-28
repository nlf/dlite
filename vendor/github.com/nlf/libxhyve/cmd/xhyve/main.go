package main

import (
	"fmt"
	"os"

	"github.com/nlf/libxhyve"
)


func main() {
	done := make(chan bool)
	ptyCh := make(chan string)

	go func() {
		if err := xhyve.Run(os.Args, ptyCh); err != nil {
			fmt.Println(err)
		}
		done <- true
	}()

	if len(os.Args) <= 1 {
		fmt.Println("No arguments found, there is nothing to do.")
		return
	}

	fmt.Printf("Waiting on a pseudo-terminal to be ready... ")
	pty := <-ptyCh
	fmt.Println("done")
	fmt.Printf("PTY allocated at %s\n", pty)

	<-done
	fmt.Println("Hypervisor goroutine finished!")
}
