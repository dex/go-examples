package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Need interface name")
	}

	iface, err := net.InterfaceByName(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(fd)

	buf := []byte{0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x55}
	err = syscall.Sendto(fd, buf, 0, &syscall.SockaddrLinklayer{Ifindex: iface.Index})
	if err != nil {
		log.Fatal("Cannot sendto:", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
