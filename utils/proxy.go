package utils

import (
	"io"
	"net"
	"os"
	"sync"
)

func Proxy(ip string) error {
	_, err := os.Stat("/var/run/docker.sock")
	if err == nil {
		err = os.Remove("/var/run/docker.sock")
		if err != nil {
			return err
		}
	}

	listener, err := net.Listen("unix", "/var/run/docker.sock")
	if err != nil {
		return err
	}

	err = os.Chmod("/var/run/docker.sock", 0777)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for {
		client, err := listener.Accept()
		if err != nil {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			// defer client.Close()
			server, err := net.Dial("tcp", ip+":2375")
			if err != nil {
				return
			}
			// defer server.Close()
			go func() {
				_, err := io.Copy(client, server)
				if err != nil {
					return
				}

				client.Close()
				server.Close()
			}()

			go func () {
				_, err := io.Copy(server, client)
				if err != nil {
					return
				}

				client.Close()
				server.Close()
			}()
		}()
	}

	wg.Wait()
	return nil
}
