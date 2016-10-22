package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
)

type API struct {
	done   chan bool
	daemon *Daemon
}

func extractUser(r *http.Request) (*User, error) {
	userHeader, ok := r.Header["X-Username"]
	if !ok {
		return nil, fmt.Errorf("Missing X-Username header")
	}

	return lookupUser(userHeader[0])
}

func (a *API) start(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	if a.daemon.VM != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Conflict"))
		return
	}

	user, err := extractUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	a.daemon.VM, err = NewVM(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = a.daemon.VM.Start()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(fmt.Sprintf("Virtual machine started, tty available at %s", a.daemon.VM.TTY)))
}

func (a *API) started(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	if a.daemon.VM == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Virtual machine has not been started"))
		return
	}

	a.daemon.VM.Started = true
	w.Write([]byte("Virtual machine flagged as started"))
}

func (a *API) stop(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	if a.daemon.VM == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Virtual machine is not running"))
		return
	}

	a.daemon.VM.Started = false
	a.daemon.VM.Ready = false
	a.daemon.VM.Stop()
	a.daemon.VM = nil

	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Virtual machine shut down"))
}

func (a *API) status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if a.daemon.VM == nil {
		user, err := extractUser(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("\"message\": \"Unauthorized\""))
			return
		}

		status, err := EmptyStatus(*user)
		if err != nil {
			statusErr := VMStatusError{
				Status:  "error",
				Message: err.Error(),
			}
			js, _ := json.Marshal(statusErr)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(js)
			return
		}

		js, _ := json.Marshal(status)

		w.WriteHeader(http.StatusOK)
		w.Write(js)
		return
	}

	status, err := a.daemon.VM.Status()
	if err != nil {
		statusErr := VMStatusError{
			Status:  "error",
			Message: err.Error(),
		}
		js, _ := json.Marshal(statusErr)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(js)
		return
	}

	js, err := json.Marshal(status)
	if err != nil {
		statusErr := VMStatusError{
			Status:  "error",
			Message: err.Error(),
		}
		js, _ := json.Marshal(statusErr)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(js)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (a *API) Listen() error {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:1050")
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
	mux.HandleFunc("/start", a.start)
	mux.HandleFunc("/started", a.started)
	mux.HandleFunc("/stop", a.stop)
	mux.HandleFunc("/status", a.status)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	loggedMux := newLoggedHandler(mux, os.Stdout)
	server := &http.Server{
		Handler: loggedMux,
	}

	return server.Serve(listener)
}

func (a *API) Stop() {
	a.done <- true
}

func NewAPI(daemon *Daemon) *API {
	return &API{
		daemon: daemon,
		done:   make(chan bool),
	}
}
