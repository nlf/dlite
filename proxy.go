package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

type Proxy struct {
	done   chan bool
	daemon *Daemon
}

func (p *Proxy) cleanup() error {
	_, err := os.Stat("/var/run/docker.sock")
	if err == nil {
		err = os.Remove("/var/run/docker.sock")
		return err
	}

	return nil
}

func (p *Proxy) proxy(w http.ResponseWriter, r *http.Request) {
	if p.daemon.VM == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("The virtual machine has not been started"))
		return
	}

	addr, err := p.daemon.VM.Address()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Unable to locate the virtual machine"))
		return
	}

	backend, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Unable to connect to the virtual machine"))
		return
	}
	defer backend.Close()

	r.URL.Scheme = "http"
	r.URL.Host = fmt.Sprintf("%s:%d", addr.IP.String(), addr.Port)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Unable to create hijacker"))
		return
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Unable to hijack connection"))
		return
	}

	r.Write(backend)
	finished := make(chan error, 1)

	go func(backend *net.TCPConn, conn net.Conn, finished chan error) {
		buf := make([]byte, 8092)
		_, err := io.CopyBuffer(backend, conn, buf)
		backend.CloseWrite()
		finished <- err
	}(backend, conn, finished)

	go func(backend *net.TCPConn, conn net.Conn, finished chan error) {
		buf := make([]byte, 8092)
		_, err := io.CopyBuffer(conn, backend, buf)
		conn.Close()
		finished <- err
	}(backend, conn, finished)

	<-finished
	<-finished
}

func (p *Proxy) Listen() error {
	err := p.cleanup()
	if err != nil {
		return err
	}

	addr, err := net.ResolveUnixAddr("unix", "/var/run/docker.sock")
	if err != nil {
		return err
	}

	raw, err := net.ListenUnix("unix", addr)
	if err != nil {
		return err
	}

	listener := &unixListener{
		UnixListener: raw,
		done:         p.done,
	}

	err = os.Chmod("/var/run/docker.sock", 0777)
	if err != nil {
		return err
	}

	server := http.Server{
		Handler: http.HandlerFunc(p.proxy),
	}

	return server.Serve(listener)
}

func (p *Proxy) Stop() {
	p.done <- true
	p.cleanup()
}

func NewProxy(daemon *Daemon) *Proxy {
	return &Proxy{
		daemon: daemon,
		done:   make(chan bool),
	}
}
