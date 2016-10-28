# DLite

this branch represents the latest *beta* version. the stable version can be found in the legacy branch.

## Building

install dependencies

```sh
brew install opam golang
opam init
eval `opam config env`
opam pin add qcow-format git://github.com/mirage/ocaml-qcow#master
opam install uri qcow-format
go get -u github.com/jteeuwen/go-bindata/...
git submodule init
git submodule update
```

update dependencies (use this if you've already built the project before)

```sh
git submodule foreach git pull origin master
opam update
opam upgrade
```

build the binary

```sh
go generate
go build
```

## TODO

- write uninstall command to remove daemon plist, resolver config, nfs exports, ssh config, and user's instance
- investigate what we need to support vpn users
