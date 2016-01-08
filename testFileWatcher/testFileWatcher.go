package main

import (
	"gopkg.in/fsnotify.v1"
	"log"
	"os"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	if err := watcher.Add(os.Args[1]); err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case event := <-watcher.Events:
			log.Print("event: ", event)
		case err := <-watcher.Errors:
			log.Print("error: ", err)
		}
	}
}
