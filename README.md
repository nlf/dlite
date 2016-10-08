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

- refactor `vm.go` to shell out to com.docker.hyperkit and track pid, figure out what state information we want to store (uuid, boot status, docker status, docker version, etc)
- write `net.go` to find proper subnet/ip for host, and lookup ip for vm. also figure out dns/hosts file/sshconfig (which is best?)
- write `nfs.go` to create exports (we're running away from 9p) and make sure `nfsd` is running, call it from `cmd_init.go`
- write `ssh.go` to create keypair for vm
- add route to `api.go` to allow vm to download a config tarball on first boot instead of passing everything on the kernel cmdline
- write `cmd_start.go` and `cmd_stop.go`, simple HTTP calls to the api
- write `cmd_status.go` to display vm's status, add option to allow printing only ip
- write `cmd_ssh.go` as a shortcut to ssh to vm
- create dlite-os (fork from dhyve-os), modify to fetch config from api on boot. add init scripts to inform api of vm's status (i.e. network is up, docker is running, etc)
- figure out how to integrate with brew services
