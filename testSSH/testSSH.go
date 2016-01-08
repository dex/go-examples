// testSSH is an example of SSH client
package main

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// read key file into buffer
	privKey, err := ioutil.ReadFile("/home/dex/.ssh/id_rsa")
	if err != nil {
		log.Fatal(err)
	}

	// parse private key in buffer
	signer, err := ssh.ParsePrivateKey(privKey)
	if err != nil {
		log.Fatal(err)
	}

	config := &ssh.ClientConfig{
		User: "dex",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", "kvm.dex.idv.tw:22", config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// request pty
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatal(err)
	}
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr
	session.Stdout = os.Stdout

	// configure local terminal via fd 0 (stdin)
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatal(err)
	}
	defer terminal.Restore(0, oldState)

	// run shell
	if err := session.Shell(); err != nil {
		log.Fatal(err)
	}

	// wait for remote shell exit
	if err := session.Wait(); err != nil {
		log.Fatal(err)
	}
}
