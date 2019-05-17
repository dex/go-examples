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
	"strings"

	"github.com/dex/iso9660"
)

var port = flag.Int("p", 8080, "`Port` to linsten on")
var mac = flag.String("m", "", "client `mac` address")

type isoFs struct {
	fs *iso9660.FileSystem
}

func (iso *isoFs) Open(name string) (http.File, error) {
	return iso.fs.Open(name)
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		os.Exit(1)
	}

	fs, err := iso9660.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/repo/", http.StripPrefix("/repo/", http.FileServer(&isoFs{fs})))
	http.HandleFunc("/v1/boot/", api)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func api(w http.ResponseWriter, r *http.Request) {
	if len(*mac) > 0 && strings.Compare(strings.ToLower(*mac), filepath.Base(r.URL.Path)) != 0 {
		log.Printf("Ignore request from %s", filepath.Base(r.URL.Path))
		http.NotFound(w, r)
		return
	}
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
		C: fmt.Sprintf("ip=dhcp inst.repo=http://%s/repo inst.ks=http://%s/repo/install/install.ks console=ttyS0,115200", r.Host, r.Host),
	}

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		panic(err)
	}
}
