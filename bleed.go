package main

import (
	bleed "github.com/FiloSottile/Heartbleed/bleed"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

var usageMessage = `This is a tool for detecting OpenSSL Heartbleed vulnerability (CVE-2014-0160).

Usage:  %s [flags] server_name[:port]

The default port is 443 (HTTPS).

The following flags are recognized:
`

func usage() {
	fmt.Fprintf(os.Stderr, usageMessage, os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	var tgt bleed.Target

	flag.StringVar(&tgt.Service, "service", "https", fmt.Sprintf("Specify a service name to test (using STARTTLS if necessary). \n\t\tBesides HTTPS, currently supported services are: \n\t\t%s", bleed.Services))
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
	}

	tgt.HostIp = flag.Arg(0)

	u, err := url.Parse(tgt.HostIp)
	if err == nil && u.Host != "" {
		tgt.HostIp = u.Host
	}

	out, err := bleed.Heartbleed(&tgt, []byte("heartbleed.filippo.io"))
	if err == bleed.Safe {
		log.Printf("%v - SAFE", tgt.HostIp)
		os.Exit(0)
	} else if err != nil {
		log.Printf("%v - ERROR: %v", tgt.HostIp, err)
		os.Exit(2)
	} else {
		log.Printf("%v\n", string(out))
		log.Printf("%v - VULNERABLE", tgt.HostIp)
		os.Exit(1)
	}
}
