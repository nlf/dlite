package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LoggedRequest struct {
	http.ResponseWriter

	req    *http.Request
	output io.Writer
	code   int
	start  time.Time
	finish time.Time
}

func (l *LoggedRequest) WriteHeader(code int) {
	l.code = code
	l.ResponseWriter.WriteHeader(code)
}

func (l *LoggedRequest) printLog() {
	ip := l.req.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon > -1 {
		ip = ip[:colon]
	}

	elapsed := l.finish.Sub(l.start).Nanoseconds() / 1000
	finish := l.finish.Format("010206.150405")
	fmt.Fprintf(l.output, "%s %s %s %s [%d] %dms\n", ip, finish, l.req.Method, l.req.RequestURI, l.code, elapsed)
}

type LoggingHandler struct {
	handler http.Handler
	output  io.Writer
}

func (h *LoggingHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	entry := &LoggedRequest{
		ResponseWriter: rw,
		output:         h.output,
		req:            r,
		start:          time.Now(),
	}

	entry.start = time.Now()
	h.handler.ServeHTTP(entry, r)
	entry.finish = time.Now()
	entry.printLog()
}

func newLoggedHandler(handler http.Handler, output io.Writer) http.Handler {
	return &LoggingHandler{
		handler: handler,
		output:  output,
	}
}
