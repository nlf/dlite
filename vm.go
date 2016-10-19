package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
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

	return os.Chown(filepath.Join(vm.Path, "vm.tty"), vm.Owner.Uid, vm.Owner.Gid)
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
			"-s", fmt.Sprintf("4:0,virtio-blk,file://%s,direct,format=qcow", disk),
			"-f", fmt.Sprintf("kexec,%s,%s,earlyprintk=serial console=ttyS0 config=%s", kernel, rootfs, bootstrapData),
		},
	}, nil
}
