package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	remote    = "192.168.2.112:23233"
	cert_path = "/tftpboot/cert.pem"
	key_path  = "/tftpboot/key.pem"
)

func main() {
	buf := make([]byte, 1024)
	cert, err := tls.LoadX509KeyPair(cert_path, key_path)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := tls.Dial("tcp", remote, &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Local address
	fmt.Printf("%s\n", conn.LocalAddr().(*net.TCPAddr).String())

	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	buf = buf[:n]
	s := string(buf)
	list := strings.Split(s, ",")
	for _, t := range list {
		fmt.Println(strings.Trim(t, " \n\t"))
	}
	fmt.Println(n)
	fmt.Println(conn.ConnectionState())
}
