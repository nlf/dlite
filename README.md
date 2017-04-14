# DLite

this branch represents the latest *beta* version. the stable version can be found in the legacy branch.

## Building

install dependencies

```sh
brew install opam golang libev
opam init
eval `opam config env`
opam install uri qcow.0.7.0 conf-libev logs fmt qcow-format
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

## Binary

To use the binary, you'll need `homebrew` to install `libev`.

```sh
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
brew install libev
```
