DEPS := $(shell find . -name '*.go')

all: dlite dlitesvc

dlite: ${DEPS}
	go build -ldflags="-s -w"

dlitesvc: ${DEPS}
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dlitesvc

clean:
	go clean
	rm -f dlite dlitesvc

.phony: all clean
