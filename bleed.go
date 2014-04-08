package main

import (
	"flag"
	"fmt"
	bleed "github.com/FiloSottile/Heartbleed/bleed"
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
	var tgt bleed.Target

	flag.StringVar(&tgt.StartTls, "starttls", "", "use STARTTLS")
	flag.Parse()

	if flag.NArg() < 1 {
		usage(os.Args[0])
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
