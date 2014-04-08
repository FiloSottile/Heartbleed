package main

import (
	bleed "github.com/FiloSottile/Heartbleed/bleed"
	"fmt"
	"log"
	"net/url"
	"os"
)

var usageMessage = `This is a tool for detecting OpenSSL Heartbleed vulnerability (CVE-2014-0160).

Usage:

	%s server_name(:port)
	
	The default port is 443 (HTTPS).
`

func usage(progname string) {
	fmt.Fprintf(os.Stderr, usageMessage, progname)
	os.Exit(2)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		usage(args[0])
	}

	host := args[1]
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
