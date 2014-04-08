package main

import (
	"flag"
	bleed "github.com/FiloSottile/Heartbleed/bleed"
	"log"
	"os"
)

func main() {
	var tgt bleed.Target

	flag.StringVar(&tgt.StartTls, "starttls", "", "use STARTTLS")
	flag.Parse()

	tgt.HostIp = flag.Arg(0)

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
