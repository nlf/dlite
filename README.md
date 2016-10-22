# DLite

## Building

install dependencies

```sh
brew install opam golang
opam init
eval `opam config env`
opam pin add qcow-format git://github.com/mirage/ocaml-qcow#master
opam install uri qcow-format
go get -u github.com/jteeuwen/go-bindata/...
```

update dependencies (use this if you've already built the project before)

```sh
git submodule update --init
opam update
opam upgrade
```

build the binary

```sh
go generate
go build
```

## TODO

- write `net.go` to find proper subnet/ip for host, and lookup ip for vm. also figure out dns/hosts file/sshconfig (which is best?)
- write `ssh.go` to create keypair for vm
- write `cmd_ssh.go` as a shortcut to ssh to vm
- write `cmd_tty.go` as a shortcut to open a screen terminal to the vm
- figure out how to integrate with brew services
