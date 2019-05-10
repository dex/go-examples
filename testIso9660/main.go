// Package main provides ...
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dex/iso9660"
)

type isoFs struct {
	fs *iso9660.FileSystem
}

func (iso *isoFs) Open(name string) (http.File, error) {
	return iso.fs.Open(name)
}

func main() {
	argc := len(os.Args)

	if argc < 2 {
		os.Exit(1)
	}

	fs, err := iso9660.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", http.FileServer(&isoFs{fs})))

}
