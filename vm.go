package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var (
	blockMatcher = regexp.MustCompile(`\{((?s)[^\{]*)\}`)
	nameMatcher  = regexp.MustCompile(`(?s).*name=([a-fA-F0-9\-]+)`)
	ipMatcher    = regexp.MustCompile(`(?s).*ip_address=([0-9\.]+)`)
)

type dockerNetwork struct {
	Driver string
	IPAM   struct {
		Config []struct {
			Subnet string
		}
	}
}

type containerNetwork struct {
	NetworkSettings struct {
		IPAddress string
	}
}

type VMStatusError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type VMStatus struct {
	Config

	Started bool   `json:"started"`
	IP      string `json:"ip"`
	Pid     int    `json:"pid"`
}

type VM struct {
	Config  Config
	Owner   *User
	TTY     string
	Log     string
	Path    string
	Bin     string
	Args    []string
	Process *exec.Cmd
	Started bool
	Ready   bool
	ip      string
}

func (vm *VM) Start() error {
	if vm.Process != nil {
		return fmt.Errorf("Virtual machine is already running")
	}

	vm.Process = exec.Command(vm.Bin, vm.Args...)
	err := vm.Process.Start()
	if err != nil {
		return err
	}

	err = vm.waitForBoot()
	if err != nil {
		return err
	}

	err = os.Chown(filepath.Join(vm.Path, "vm.tty"), vm.Owner.Uid, vm.Owner.Gid)
	if err != nil {
		return err
	}

	if vm.Config.Route {
		err = vm.Route()
		if err != nil {
			return err
		}
	}

	return nil
}

func (vm *VM) waitForBoot() error {
	attempts := 0
	success := fmt.Sprintf("%s login:", vm.Config.Hostname)

	for {
		if attempts >= 30 {
			return fmt.Errorf("Timed out waiting for virtual machine")
		}

		time.Sleep(time.Second)

		file, err := os.Open(vm.Log)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), success) {
				return nil
			}
		}

		attempts++
	}
}

func (vm *VM) Route() error {
	subnet, err := vm.dockerSubnet()
	if err != nil {
		return err
	}

	err = exec.Command("route", "-n", "add", subnet, vm.Config.Hostname).Run()
	if err != nil {
		return err
	}

	routeBytes, err := exec.Command("route", "-n", "get", vm.Config.Hostname).Output()
	if err != nil {
		return err
	}

	routeIfaceMatcher := regexp.MustCompile(`(?m)^\s*interface:\s*(\w+)$`)
	routeIfaceMatches := routeIfaceMatcher.FindAllStringSubmatch(string(routeBytes), -1)
	if routeIfaceMatches == nil {
		return fmt.Errorf("Unable to find interface")
	}

	routeIface := routeIfaceMatches[0][1]

	memberBytes, err := exec.Command("ifconfig", routeIface).Output()
	if err != nil {
		return err
	}

	memberMatcher := regexp.MustCompile(`(?m)^\s*member:\s*(.*) flags.*$`)
	memberMatches := memberMatcher.FindAllStringSubmatch(string(memberBytes), -1)
	if memberMatches == nil {
		return fmt.Errorf("Unable to find interface members")
	}

	members := strings.Split(memberMatches[0][1], " ")
	for _, member := range members {
		err := exec.Command("ifconfig", routeIface, "-hostfilter", member).Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func (vm *VM) dockerSubnet() (string, error) {
	ip, err := vm.IP()
	if err != nil {
		return "", err
	}

	res, err := http.Get(fmt.Sprintf("http://%s:2375/networks?filter={\"type\":{\"builtin\":true}}", ip))
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	networks := []dockerNetwork{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&networks)
	if err != nil {
		return "", err
	}

	for _, network := range networks {
		if network.Driver == "bridge" {
			return network.IPAM.Config[0].Subnet, nil
		}
	}

	return "", fmt.Errorf("Unable to find bridge network")
}

func (vm *VM) findContainer(name string) (string, error) {
	ip, err := vm.IP()
	if err != nil {
		return "", err
	}

	res, err := http.Get(fmt.Sprintf("http://%s:2375/containers/%s/json", ip, name))
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	container := containerNetwork{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&container)
	if err != nil {
		return "", err
	}

	if container.NetworkSettings.IPAddress != "" {
		return container.NetworkSettings.IPAddress, nil
	}

	return "", fmt.Errorf("Unable to find container")
}

func (vm *VM) Stop() error {
	if vm.Process == nil {
		return nil
	}

	err := vm.Process.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}

	err = vm.Process.Wait()
	vm.Process = nil
	os.RemoveAll(filepath.Join(vm.Path, "vm.tty"))
	return err
}

func (vm *VM) IP() (string, error) {
	if vm.ip != "" {
		return vm.ip, nil
	}

	type result struct {
		value string
		err   error
	}

	value := make(chan result, 1)

	go func() {
		attempts := 0
		for {
			if attempts >= 15 {
				value <- result{"", fmt.Errorf("Timed out waiting for IP address")}
				break
			}

			time.Sleep(time.Second)

			file, err := os.Open("/var/db/dhcpd_leases")
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}

				value <- result{"", err}
				break
			}

			defer file.Close()
			leases, err := ioutil.ReadAll(file)
			if err != nil {
				value <- result{"", err}
				break
			}

			blocks := blockMatcher.FindAllStringSubmatch(string(leases), -1)
			for _, block := range blocks {
				name := nameMatcher.FindStringSubmatch(block[1])
				if name != nil && name[1] == vm.Config.Id {
					ip := ipMatcher.FindStringSubmatch(block[1])
					if ip != nil {
						value <- result{ip[1], nil}
						break
					}
				}
			}

			attempts++
		}
	}()

	res := <-value
	if res.err == nil {
		vm.ip = res.value
	}

	return res.value, res.err
}

func (vm *VM) Address() (*net.TCPAddr, error) {
	ip, err := vm.IP()
	if err != nil {
		return nil, err
	}

	return net.ResolveTCPAddr("tcp", ip+":2375")
}

func (vm *VM) Status() (VMStatus, error) {
	status := VMStatus{}
	cfg, err := readConfig(getPath(*vm.Owner))
	if err != nil {
		return status, err
	}

	status.Config = cfg

	if vm.Process != nil {
		ip, err := vm.IP()
		if err != nil {
			return status, nil
		}

		status.IP = ip
		status.Started = true
		status.Pid = vm.Process.Process.Pid
	}

	return status, nil
}

func EmptyStatus(owner User) (VMStatus, error) {
	status := VMStatus{}
	cfg, err := readConfig(getPath(owner))
	if err != nil {
		return status, err
	}

	status.Config = cfg
	return status, nil
}

func NewVM(owner *User) (*VM, error) {
	path := getPath(*owner)
	cfg, err := readConfig(path)
	if err != nil {
		return nil, err
	}

	kernel := filepath.Join(path, "bzImage")
	rootfs := filepath.Join(path, "rootfs.cpio.xz")
	disk := filepath.Join(path, "disk.qcow")
	tty := filepath.Join(path, "vm.tty")
	log := filepath.Join(path, "vm.log")
	bootstrapData, err := GetBootstrapData(*owner)
	if err != nil {
		return nil, err
	}

	return &VM{
		Owner:  owner,
		Config: cfg,
		Path:   path,
		Bin:    filepath.Join(path, "bin", "com.docker.hyperkit"),
		TTY:    tty,
		Log:    log,
		Args: []string{
			"-A",
			"-u",
			"-U", cfg.Id,
			"-c", fmt.Sprintf("%d", cfg.Cpu),
			"-m", fmt.Sprintf("%dG", cfg.Memory),
			"-l", fmt.Sprintf("com1,autopty=%s,log=%s", tty, log),
			"-s", "0:0,hostbridge",
			"-s", "5,virtio-rnd",
			"-s", "31,lpc",
			"-s", "2:0,virtio-net",
			"-s", fmt.Sprintf("4:0,virtio-blk,file://%s,format=qcow", disk),
			"-f", fmt.Sprintf("kexec,%s,%s,earlyprintk=serial console=ttyS0 no_timer_check config=%s", kernel, rootfs, bootstrapData),
		},
	}, nil
}
