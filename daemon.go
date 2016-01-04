package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/TheNewNormal/libxhyve"
	"github.com/codegangsta/cli"
)

func daemonHandler(ctx *cli.Context) {
	config := readConfig()
	go proxySocket(config.Uuid)
	runXhyve(config)
}

func runXhyve(config Config) {
	done := make(chan bool)
	pty := make(chan string)
	args := []string{
		"libxhyve_bug",
		"-A",
		"-c", fmt.Sprintf("%d", config.CpuCount),
		"-m", fmt.Sprintf("%dG", config.Memory),
		"-s", "0:0,hostbridge",
		"-l", "com1,stdio",
		"-s", "31,lpc",
		"-s", "2:0,virtio-net",
		"-s", "4,virtio-blk,"+PREFIX+"/var/db/dlite/disk.img",
		"-U", config.Uuid,
		"-f", fmt.Sprintf("kexec,%s,%s,%q", PREFIX+"/usr/share/dlite/bzImage", PREFIX+"/usr/share/dlite/rootfs.cpio.xz", "console=ttyS0 hostname=dlite uuid="+config.Uuid),
	}

	go func(args []string, pty chan string) {
		err := xhyve.Run(args, pty)
		if err != nil {
			log.Fatalln(err)
		}

		done <- true
	}(args, pty)

	ptyStr := <-pty
	log.Printf("Xhyve listening on term: %s\n", ptyStr)
	<-done
}

func proxySocket(uuid string) {
	ip := getIP(uuid)
	log.Println("got ip:", ip)
	for ip == "" {
		ip = getIP(uuid)
		log.Println(ip)
	}

	_, err := os.Stat("/var/run/docker.sock")
	if err == nil {
		err = os.Remove("/var/run/docker.sock")
		if err != nil {
			log.Fatal(err)
		}
	}
	socket, err := net.Listen("unix", "/var/run/docker.sock")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chmod("/var/run/docker.sock", 0777)
	if err != nil {
		log.Fatal(err)
	}

	for {
		log.Print("Waiting for connection to socket...")
		conn, err := socket.Accept()
		log.Print("Accepted connection from socket..")
		if err != nil {
			log.Fatal(err)
		}

		docker, err := net.Dial("tcp", ip+":2375")
		if err != nil {
			log.Fatal(err)
		}

		defer docker.Close()

		go func(c, d net.Conn) {
			_, err = io.Copy(d, c)
			if err != nil {
				log.Fatal(err)
			}
		}(conn, docker)

		go func(c, d net.Conn) {
			_, err = io.Copy(c, d)
			if err != nil {
				log.Fatal(err)
			}
		}(conn, docker)
	}
}
