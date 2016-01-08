package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("ip:gre", "10.42.0.3")
	if err != nil {
		log.Fatal("Cannot create raw socket for GRE")
	}
	buf := make([]byte, 1518)
	_, err = conn.Write([]byte{0x0, 0x0, 0x65, 0x58, 0xb8, 0xa3, 0x86, 0x6f, 0x15, 0x86, 0xb8, 0xa3, 0x86, 0x6f, 0x15, 0x87, 0x08, 0x00})
	if err != nil {
		log.Fatal("Cannot write raw socket for GRE")
	}
	for {
		// read data from raw socket
		numRead, addr, _ := conn.(*net.IPConn).ReadFrom(buf)
		fmt.Printf("Read %d bytes from %s\n", numRead, addr)
	}
}
