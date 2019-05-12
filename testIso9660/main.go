// Package main provides ...
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dex/iso9660"
)

var port = flag.Int("port", 8080, "Port to linsten on")

type isoFs struct {
	fs *iso9660.FileSystem
}

func (iso *isoFs) Open(name string) (http.File, error) {
	return iso.fs.Open(name)
}

func main() {
	flag.Parse()
	argc := len(os.Args)

	if argc < 2 {
		os.Exit(1)
	}

	fs, err := iso9660.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/repo/", http.StripPrefix("/repo/", http.FileServer(&isoFs{fs})))
	http.HandleFunc("/v1/boot/", api)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))

}

func api(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving boot config for %s", filepath.Base(r.URL.Path))
	resp := struct {
		K string   `json:"kernel"`
		I []string `json:"initrd"`
		C string   `json:"cmdline"`
	}{
		K: "/repo/isolinux/vmlinuz",
		I: []string{
			"/repo/isolinux/initrd.img",
		},
		C: fmt.Sprintf("ip=dhcp inst.repo=http://%s/repo inst.ks=http://10.206.83.38:8000/install.ks", r.Host),
	}

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		panic(err)
	}
}
