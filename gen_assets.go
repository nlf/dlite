package main

// build hyperkit
//go:generate make -C hyperkit

// copy binary assets
//go:generate sh -c "cp `which qcow-tool` assets/"
//go:generate cp hyperkit/build/com.docker.hyperkit assets/

// generate bundled assets
//go:generate go-bindata -pkg main -o assets.go -prefix assets assets/
