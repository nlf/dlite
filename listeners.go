package main

import (
	"fmt"
	"net"
	"time"
)

type tcpListener struct {
	*net.TCPListener
	done chan bool
}

func (l *tcpListener) Accept() (net.Conn, error) {
	for {
		l.SetDeadline(time.Now().Add(time.Second))
		select {
		case <-l.done:
			return nil, fmt.Errorf("Server closed")
		default:
		}

		newConn, err := l.TCPListener.Accept()
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, err
	}
}

func (l *tcpListener) Close() error {
	close(l.done)
	return nil
}

type unixListener struct {
	*net.UnixListener
	done chan bool
}

func (s *unixListener) Accept() (net.Conn, error) {
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

func (s *unixListener) Close() error {
	close(s.done)
	return nil
}
