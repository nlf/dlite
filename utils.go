package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func getIP(uuid string) string {
	nameRe := regexp.MustCompile(`.*name=([A-F0-9\-]+)`)
	ipRe := regexp.MustCompile(`.*ip_address=([0-9\.]+)`)

	file, err := os.Open("/var/db/dhcpd_leases")
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()
	leases, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	lines := string(leases)
	names := nameRe.FindStringSubmatch(lines)
	ips := ipRe.FindStringSubmatch(lines)
	for i, name := range names {
		log.Println("name", name)
		if name == uuid {
			return ips[i]
		}
	}

	return ""
}
