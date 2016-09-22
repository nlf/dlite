package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

type stoppableListener struct {
	*net.UnixListener
	done chan bool
}

func (s *stoppableListener) Accept() (net.Conn, error) {
	for {
		s.SetDeadline(time.Now().Add(time.Second))
		select {
		case <-s.done:
			return nil, fmt.Errorf("Server closed")
		default:
		}

		newConn, err := s.UnixListener.Accept()
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, err
	}
}

func (s *stoppableListener) Close() error {
	close(s.done)
	return nil
}

type Proxy struct {
	done      chan bool
	vm        *VM
	vmAddress *net.TCPAddr
}

func (p *Proxy) cleanup() error {
	_, err := os.Stat("/var/run/docker.sock")
	if err == nil {
		err = os.Remove("/var/run/docker.sock")
		return err
	}

	return nil
}

func (p *Proxy) handler(w http.ResponseWriter, r *http.Request) {
	if p.vmAddress == nil {
		addr, err := p.vm.IP()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Unable to locate the virtual machine"))
			return
		}

		p.vmAddress = addr
	}

	backend, err := net.DialTCP("tcp", nil, p.vmAddress)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Unable to connect to the virtual machine"))
		return
	}
	defer backend.Close()

	r.URL.Scheme = "http"
	r.URL.Host = fmt.Sprintf("%s:%d", p.vmAddress.IP.String(), p.vmAddress.Port)

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

	listener := &stoppableListener{
		UnixListener: raw,
		done:         p.done,
	}

	err = os.Chmod("/var/run/docker.sock", 0777)
	if err != nil {
		return err
	}

	server := http.Server{
		Handler: http.HandlerFunc(p.handler),
	}

	return server.Serve(listener)
}

func (p *Proxy) Stop() {
	p.done <- true
}

func NewProxy() (*Proxy, error) {
	vm, err := NewVM()
	if err != nil {
		return nil, err
	}

	return &Proxy{
		done: make(chan bool),
		vm:   vm,
	}, nil
}
