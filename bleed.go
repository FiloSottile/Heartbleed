package main

import (
	bleed "github.com/FiloSottile/Heartbleed/bleed"
	"log"
	"net/url"
	"os"
)

func main() {
	host := os.Args[1]
	u, err := url.Parse(host)
	if err == nil && u.Host != "" {
		host = u.Host
	}
	out, err := bleed.Heartbleed(host, []byte("heartbleed.filippo.io"))
	if err == bleed.Safe {
		log.Printf("%v - SAFE", host)
		os.Exit(0)
	} else if err != nil {
		log.Printf("%v - ERROR: %v", host, err)
		os.Exit(2)
	} else {
		log.Printf("%v\n", string(out))
		log.Printf("%v - VULNERABLE", host)
		os.Exit(1)
	}
}
