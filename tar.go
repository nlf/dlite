package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func generateTarball(user User) ([]byte, error) {
	basePath := getPath(user)
	cfg, err := readConfig(basePath)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	gz := gzip.NewWriter(buf)
	tarball := tar.NewWriter(gz)

	hostname := []byte(cfg.Hostname)
	hostnameHeader := &tar.Header{
		Name:  "/etc/hostname",
		Mode:  0644,
		Size:  int64(len(hostname)),
		Uname: "root",
		Gname: "wheel",
	}

	err = tarball.WriteHeader(hostnameHeader)
	if err != nil {
		return nil, err
	}

	_, err = tarball.Write(hostname)
	if err != nil {
		return nil, err
	}

	hosts := []byte(fmt.Sprintf("127.0.0.1 localhost %s", cfg.Hostname))
	hostsHeader := &tar.Header{
		Name:  "/etc/hosts",
		Mode:  0644,
		Size:  int64(len(hosts)),
		Uname: "root",
		Gname: "wheel",
	}

	err = tarball.WriteHeader(hostsHeader)
	if err != nil {
		return nil, err
	}

	_, err = tarball.Write(hosts)
	if err != nil {
		return nil, err
	}

	ifaces := []byte(fmt.Sprintf("auto lo\niface lo inet loopback\n\nauto eth0\niface eth0 inet dhcp\nhostname %s", cfg.Id))
	ifacesHeader := &tar.Header{
		Name:  "/etc/network/interfaces",
		Mode:  0644,
		Size:  int64(len(ifaces)),
		Uname: "root",
		Gname: "wheel",
	}

	err = tarball.WriteHeader(ifacesHeader)
	if err != nil {
		return nil, err
	}

	_, err = tarball.Write(ifaces)
	if err != nil {
		return nil, err
	}

	dns := []byte(fmt.Sprintf("nameserver %s", cfg.DNS))
	dnsHeader := &tar.Header{
		Name:  "/etc/resolv.conf",
		Mode:  0644,
		Size:  int64(len(dns)),
		Uname: "root",
		Gname: "wheel",
	}

	err = tarball.WriteHeader(dnsHeader)
	if err != nil {
		return nil, err
	}

	_, err = tarball.Write(dns)
	if err != nil {
		return nil, err
	}

	hostIp, err := getHostAddress()
	if err != nil {
		return nil, err
	}

	hostIpBytes := []byte(hostIp)
	hostIpHeader := &tar.Header{
		Name:  "/etc/dlite/host_ip",
		Mode:  0600,
		Size:  int64(len(hostIpBytes)),
		Uname: "root",
		Gname: "wheel",
	}

	err = tarball.WriteHeader(hostIpHeader)
	if err != nil {
		return nil, err
	}

	_, err = tarball.Write(hostIpBytes)
	if err != nil {
		return nil, err
	}

	username := []byte(user.Name)
	usernameHeader := &tar.Header{
		Name:  "/etc/dlite/username",
		Mode:  0600,
		Size:  int64(len(username)),
		Uname: "root",
		Gname: "wheel",
	}

	err = tarball.WriteHeader(usernameHeader)
	if err != nil {
		return nil, err
	}

	_, err = tarball.Write(username)
	if err != nil {
		return nil, err
	}

	userId := []byte(fmt.Sprintf("%d", user.Uid))
	userIdHeader := &tar.Header{
		Name:  "/etc/dlite/userid",
		Mode:  0600,
		Size:  int64(len(userId)),
		Uname: "root",
		Gname: "wheel",
	}

	err = tarball.WriteHeader(userIdHeader)
	if err != nil {
		return nil, err
	}

	_, err = tarball.Write(userId)
	if err != nil {
		return nil, err
	}

	dockerVersion := []byte(cfg.Docker)
	dockerVersionHeader := &tar.Header{
		Name:  "/etc/dlite/docker_version",
		Mode:  0600,
		Size:  int64(len(dockerVersion)),
		Uname: "root",
		Gname: "wheel",
	}

	err = tarball.WriteHeader(dockerVersionHeader)
	if err != nil {
		return nil, err
	}

	_, err = tarball.Write(dockerVersion)
	if err != nil {
		return nil, err
	}

	dockerArgs := []byte(cfg.Extra)
	dockerArgsHeader := &tar.Header{
		Name:  "/etc/dlite/docker_args",
		Mode:  0600,
		Size:  int64(len(dockerArgs)),
		Uname: "root",
		Gname: "wheel",
	}

	err = tarball.WriteHeader(dockerArgsHeader)
	if err != nil {
		return nil, err
	}

	_, err = tarball.Write(dockerArgs)
	if err != nil {
		return nil, err
	}

	sshDirHeader := &tar.Header{
		Name:     "/home/docker/.ssh",
		Mode:     0700,
		Typeflag: tar.TypeDir,
		Uname:    "docker",
		Gname:    "docker",
	}

	err = tarball.WriteHeader(sshDirHeader)
	if err != nil {
		return nil, err
	}

	keyFile, err := os.Open(filepath.Join(basePath, "key.pub"))
	if err != nil {
		return nil, err
	}

	keyStat, err := keyFile.Stat()
	if err != nil {
		return nil, err
	}

	keysHeader := &tar.Header{
		Name:  "/home/docker/.ssh/authorized_keys",
		Mode:  0600,
		Size:  keyStat.Size(),
		Uname: "docker",
		Gname: "docker",
	}

	err = tarball.WriteHeader(keysHeader)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(tarball, keyFile)
	if err != nil {
		return nil, err
	}

	tarball.Close()
	gz.Close()

	return buf.Bytes(), nil
}

func GetBootstrapData(user User) (string, error) {
	tarball, err := generateTarball(user)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(tarball), nil
}
