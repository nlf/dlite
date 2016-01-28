DEPS := $(wildcard *.go) $(wildcard utils/*.go)

dlite: ${DEPS}
	GO15VENDOREXPERIMENT=1 go build
