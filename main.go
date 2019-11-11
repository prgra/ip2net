package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"os"
	"strings"
)

const perc = 50
const agg = 24

type elem struct {
	cnt int
	ips []net.IP
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("need file param")
		os.Exit(0)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("can't open %s: %v\n", os.Args[1], err)
	}
	defer f.Close()
	aggMask := net.CIDRMask(agg, 32)
	scanner := bufio.NewScanner(f)
	db := make(map[uint32]elem)
	for scanner.Scan() {
		sip := strings.TrimSpace(scanner.Text())
		ip := net.ParseIP(sip)
		e := db[ip2int(ip.Mask(aggMask))]
		e.cnt++
		e.ips = append(e.ips, ip)
		db[ip2int(ip.Mask(aggMask))] = e
	}
	le := math.Pow(2, 32-agg)
	for k := range db {
		if db[k].cnt > int(le)*perc/100 {
			fmt.Printf("%s/%d\n", int2ip(k), agg)
			continue
		}
		for _, ip := range db[k].ips {
			fmt.Printf("%s/32\n", ip.String())
		}
	}
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
