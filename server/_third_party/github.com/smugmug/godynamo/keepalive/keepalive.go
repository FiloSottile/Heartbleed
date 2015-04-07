package keepalive

import (
	auth_v4 "github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/auth_v4" // to get the Client
	"log"
	"time"
)

// dial the keep alive domains to establish a conn
func dialConns(keepAliveUrls []string) error {
	var e error
	for _, u := range keepAliveUrls {
		_, err := auth_v4.Client.Head(u)
		if err != nil {
			e = err
		} else {
			// log.Printf("conn %v",u) // uncomment to see keepalives
		}
	}
	return e
}

// KeepAlive can make periodic HEAD requests to our AWS endpoint url to keep conns alive.
// Should be run as a goroutine: go KeepAlive(..)
func KeepAlive(keepAliveUrls []string) {
	for {
		select {
		case <-time.After(5 * time.Second):
			dial_err := dialConns(keepAliveUrls)
			if dial_err != nil {
				log.Printf("keepalive.KeepAlive: dial fail:%s", dial_err.Error())
			}
		}
	}
}
