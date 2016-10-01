# DLite

## Building

```sh
brew install opam golang
opam init
eval `opam config env`
opam pin add qcow-format git://github.com/mirage/ocaml-qcow#master
opam install uri qcow-format
go get github.com/jteeuwen/go-bindata/...
go generate
go build
```

to update your build

```sh
git pull origin master
git submodule update --init
opam update
opam upgrade
go generate
go build
```
