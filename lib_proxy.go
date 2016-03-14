package main

import (
	// "bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

type StoppableListener struct {
	*net.UnixListener
	done chan error
}

func (s *StoppableListener) Accept() (net.Conn, error) {
	for {
		s.SetDeadline(time.Now().Add(time.Second))

		select {
		case <-s.done:
			return nil, fmt.Errorf("Xhyve exited, stopping proxy")
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

func (s *StoppableListener) Close() error {
	close(s.done)
	return nil
}

func Proxy(ip string, done chan error) error {
	_, err := os.Stat("/var/run/docker.sock")
	if err == nil {
		err = os.Remove("/var/run/docker.sock")
		if err != nil {
			return err
		}
	}

	addr, err := net.ResolveUnixAddr("unix", "/var/run/docker.sock")
	if err != nil {
		return err
	}

	rawListener, err := net.ListenUnix("unix", addr)
	if err != nil {
		return err
	}

	listener := &StoppableListener{UnixListener: rawListener, done: done}

	err = os.Chmod("/var/run/docker.sock", 0777)
	if err != nil {
		return err
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = "http"
		r.URL.Host = ip + ":2375"

		hj, ok := w.(http.Hijacker)
		if !ok {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		conn, _, err := hj.Hijack()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		addr, err := net.ResolveTCPAddr("tcp", ip+":2375")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		backend, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		r.Write(backend)

		finished := make(chan error, 1)

		go func() {
			buf := make([]byte, 8092)
			_, err := io.CopyBuffer(backend, conn, buf)
			backend.CloseWrite()
			finished <- err
		}()

		go func() {
			buf := make([]byte, 8092)
			_, err := io.CopyBuffer(conn, backend, buf)
			conn.Close()
			finished <- err
		}()

		<-finished
		<-finished
		backend.Close()
	})

	server := http.Server{
		Handler: handler,
	}

	return server.Serve(listener)
}
