package main

//go:generate sh -c "cd hyperkit; make"
//go:generate cp $HOME/.opam/system/bin/qcow-tool assets/
//go:generate cp hyperkit/build/com.docker.hyperkit assets/
//go:generate go-bindata -pkg main -o assets.go assets/
