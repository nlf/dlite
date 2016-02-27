package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

func Proxy(ip string) error {
	_, err := os.Stat("/var/run/docker.sock")
	if err == nil {
		err = os.Remove("/var/run/docker.sock")
		if err != nil {
			return err
		}
	}

	listener, err := net.Listen("unix", "/var/run/docker.sock")
	if err != nil {
		return err
	}

	err = os.Chmod("/var/run/docker.sock", 0777)
	if err != nil {
		return err
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = "http"
		r.URL.Host = ip + ":2375"

		upgrade := false
		if strings.HasSuffix(r.URL.Path, "/attach") {
			upgrade = true
		} else if len(r.Header["Upgrade"]) > 0 {
			upgrade_header := strings.ToLower(r.Header["Upgrade"][0])
			upgrade = upgrade_header == "tcp" || upgrade_header == "websocket"
		}

		if upgrade {
			hj, ok := w.(http.Hijacker)
			if !ok {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			conn, _, err := hj.Hijack()
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			defer conn.Close()

			backend, err := net.Dial("tcp", ip+":2375")
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			defer backend.Close()

			err = r.Write(backend)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			finished := make(chan bool, 1)
			go func() {
				io.Copy(backend, conn)
			}()

			go func() {
				io.Copy(conn, backend)
				finished <- true
			}()

			<-finished
		} else {
			resp, err := http.DefaultTransport.RoundTrip(r)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			defer resp.Body.Close()

			for k, v := range resp.Header {
				for _, vv := range v {
					w.Header().Add(k, vv)
				}
			}

			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		}
	})

	server := http.Server{
		Handler: handler,
	}

	return server.Serve(listener)
}
