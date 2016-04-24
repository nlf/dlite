package rpc

import (
	"net/rpc"
)

func NewClient(local bool) (*rpc.Client, error) {
	vm := new(VM)
	rpc.Register(vm)
	addr := ""
	if local {
		addr = "127.0.0.1:8899"
	} else {
		addr = "192.168.64.1:8899"
	}

	return rpc.DialHTTP("tcp", addr)
}
