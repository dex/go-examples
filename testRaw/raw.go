package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// Resolve binding address
	netaddr, err := net.ResolveIPAddr("ip6", "::")
	if err != nil {
		log.Fatal("Cannot resolve address\n")
	}
	fmt.Println("Listen on :", netaddr)

	// listen on ICMPv6
	conn, err := net.ListenIP("ip6:58", netaddr)
	if err != nil {
		log.Fatal("Cannot listen raw socket\n")
	}

	buf := make([]byte, 1024)
	for {
		// read data from raw socket
		numRead, addr, _ := conn.ReadFrom(buf)
		fmt.Printf("Read %d bytes from %s\n", numRead, addr)
	}
}
