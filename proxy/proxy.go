package proxy

import (
	"io"
	"net"
	"net/http"
	"os"

	"github.com/nlf/dlite/rpc"
)

type Proxy struct {
	server *rpc.Server
	done   chan bool
	ip     string
}

func (p *Proxy) cleanup() error {
	_, err := os.Stat("/var/run/docker.sock")
	if err == nil {
		err = os.Remove("/var/run/docker.sock")
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Proxy) handler(w http.ResponseWriter, r *http.Request) {
	if p.server.Status != "ready" {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(p.server.Status))
		return
	}

	if p.ip == "" {
		var err error
		p.ip, err = p.server.IP()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Unable to find the virtual machine's IP"))
			return
		}
	}

	addr, err := net.ResolveTCPAddr("tcp", p.ip+":2375")
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Unable to resolve TCP address"))
		return
	}

	backend, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		p.ip = ""
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Unable to connect to virtual machine"))
		return
	}
	defer backend.Close()

	r.URL.Scheme = "http"
	r.URL.Host = p.ip + ":2375"

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

func New(server *rpc.Server) *Proxy {
	return &Proxy{
		server: server,
		done:   make(chan bool),
	}
}
