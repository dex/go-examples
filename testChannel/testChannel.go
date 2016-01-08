package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	result := make(chan bool)
	go func() {
		time.Sleep(30)
		result <- true
		fmt.Println("goroutine exit")
	}()
	select {
	case <-result:
	case <-time.After(10):
		fmt.Println("main exit")
	}
	os.Exit(0)
}
