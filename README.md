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

- refactor the cli stuff because what i have now kinda sucks, need hidden commands for special processes
- add a method to allow running installation steps as root (nfs and dns setup need this)
- write a method to add ssh configuration for the vm
- write `cmd_ssh.go` as a shortcut to ssh to vm
- write `cmd_tty.go` as a shortcut to open a screen terminal to the vm
- figure out how to integrate with brew services
