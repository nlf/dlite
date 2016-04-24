package rpc

import (
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"syscall"

	"github.com/nlf/dlite/config"
	"github.com/nlf/dlite/vm"
)

type Server struct {
	cmd    *exec.Cmd
	svc    *VM
	config *config.Config
	done   chan bool
	queue  []Message
	Status string
}

func (s *Server) Listen() error {
	addr, err := net.ResolveTCPAddr("tcp", ":8899")
	if err != nil {
		return err
	}

	raw, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	listener := &stoppableListener{
		TCPListener: raw,
		done:        s.done,
	}

	finished := make(chan error)
	go func(listener *stoppableListener) {
		finished <- http.Serve(listener, nil)
	}(listener)

	<-finished
	return nil
}

func (s *Server) Stop() {
	s.cmd.Process.Signal(syscall.SIGINT)
	s.done <- true
}

func (s *Server) IP() (string, error) {
	vm := vm.New(s.config)
	return vm.IP()
}

func NewServer() (*Server, error) {
	cfg, err := config.New(os.ExpandEnv("$SUDO_USER"))
	if err != nil {
		return nil, err
	}

	err = cfg.Load()
	if err != nil {
		return nil, err
	}

	srv := &Server{
		config: cfg,
		done:   make(chan bool),
		Status: "The virtual machine has not been started",
		queue:  []Message{},
	}

	v := &VM{
		server: srv,
	}
	srv.svc = v

	rpc.Register(v)
	rpc.HandleHTTP()

	return srv, nil
}
