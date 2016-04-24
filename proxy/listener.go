package proxy

import (
	"fmt"
	"net"
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
