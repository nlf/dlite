package main

import (
	"fmt"
	"net"
	"net/http"
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

			return newConn, err
		}
	}
}

func (l *tcpListener) Close() error {
	close(l.done)
	return nil
}

type API struct {
	done   chan bool
	status string
	tty    string
}

func (a *API) start(w http.ResponseWriter, r *http.Request) {
	if a.status != "stopped" {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Virtual machine already running"))
		return
	}

	a.status = "starting"
	tty, _, err := startVM()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	a.tty = tty
	w.Write([]byte(fmt.Sprintf("Virtual machine starting, tty available at %s", a.tty)))
}

func (a *API) started(w http.ResponseWriter, r *http.Request) {
	if a.status != "starting" {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Virtual machine not in a starting state"))
		return
	}

	a.status = "started"
	w.Write([]byte("Virtual machine flagged as started"))
}

func (a *API) stop(w http.ResponseWriter, r *http.Request) {
	if a.status == "stopped" {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Virtual machine already stopped"))
		return
	}

	a.status = "stopping"
	err := stopVM()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("Virtual machine shutting down"))
}

func (a *API) stopped(w http.ResponseWriter, r *http.Request) {
	if a.status != "stopping" {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Virtual machine not in a stopping state"))
		return
	}

	a.status = "stopped"
	a.tty = ""
	w.Write([]byte("Virtual machine flagged as stopped"))
}

func (a *API) Listen() error {
	addr, err := net.ResolveTCPAddr("tcp", "192.168.64.1:1050")
	if err != nil {
		return err
	}

	raw, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	listener := &tcpListener{
		TCPListener: raw,
		done:        a.done,
	}

	mux := http.NewServeMux()
	mux.Handle("/start", http.HandlerFunc(a.start))
	mux.Handle("/started", http.HandlerFunc(a.started))
	mux.Handle("/stop", http.HandlerFunc(a.stop))
	mux.Handle("/stopped", http.HandlerFunc(a.stopped))

	server := http.Server{
		Handler: mux,
	}

	return server.Serve(listener)
}

func (a *API) Stop() {
	a.done <- true
}

func NewAPI() *API {
	return &API{
		done:   make(chan bool),
		status: "stopped",
	}
}
