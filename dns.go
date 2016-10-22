package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/miekg/dns"
)

type DNS struct {
	daemon *Daemon
	server *dns.Server
}

func (d *DNS) handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := &dns.Msg{}
	msg.SetReply(r)

	if d.daemon.VM != nil {
		msg.SetRcode(r, dns.RcodeNameError)
		w.WriteMsg(msg)
		return
	}

	userPath := getPath(*d.daemon.VM.Owner)
	cfg, err := readConfig(userPath)
	if err != nil {
		msg.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(msg)
		return
	}

	hostname := fmt.Sprintf("%s.", cfg.Hostname)
	domain := fmt.Sprintf("%s.", getDomain(cfg.Hostname))

	q := r.Question[0]
	if q.Qtype != dns.TypeA {
		msg.SetRcode(r, dns.RcodeNotImplemented)
	} else {
		if !strings.HasSuffix(q.Name, domain) {
			msg.SetRcode(r, dns.RcodeNameError)
			w.WriteMsg(msg)
			return
		}

		if q.Name == hostname {
			ip, err := d.daemon.VM.IP()
			if err != nil {
				msg.SetRcode(r, dns.RcodeNameError)
				w.WriteMsg(msg)
				return
			}

			record := &dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				A: net.ParseIP(ip).To4(),
			}

			msg.Answer = append(msg.Answer, record)
		} else {
			msg.SetRcode(r, dns.RcodeNameError)
		}
	}

	w.WriteMsg(msg)
}

func (d *DNS) Start() error {
	dns.HandleFunc(".", d.handleRequest)
	d.server = &dns.Server{
		Addr: ":1053",
		Net:  "udp",
	}

	return d.server.ListenAndServe()
}

func (d *DNS) Stop() error {
	return d.server.Shutdown()
}

func NewDNS(daemon *Daemon) *DNS {
	return &DNS{
		daemon: daemon,
	}
}

func installResolver(hostname string) error {
	err := os.MkdirAll("/etc/resolver", 0755)
	if err != nil {
		return err
	}

	domain := getDomain(hostname)
	file := fmt.Sprintf("/etc/resolver/%s", domain)
	resolver := "nameserver 127.0.0.1\nport 1053"
	return ioutil.WriteFile(file, []byte(resolver), 0644)
}
