package main

import (
	"github.com/FiloSottile/Heartbleed/bleed"
	"log"
	"os"
)

func main() {
	err := bleed.Heartbleed(os.Args[1], []byte("YELLOW SUBMARINE"))
	if err != nil {
		log.Printf("%v - SAFE", os.Args[1])
		os.Exit(1)
	} else {
		log.Printf("%v - VULNERABLE", os.Args[1])
		os.Exit(0)
	}
}
