package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/docopt/docopt-go"

	"github.com/FiloSottile/Heartbleed/heartbleed"
	"github.com/FiloSottile/Heartbleed/server/cache"
)

var PAYLOAD = []byte("filippo.io/Heartbleed")

var withCache bool

type result struct {
	Code  int    `json:"code"`
	Data  string `json:"data"`
	Error string `json:"error"`
	Host  string `json:"host"`
}

func handleRequest(tgt *heartbleed.Target, w http.ResponseWriter, r *http.Request, skip bool) {
	if tgt.HostIp == "" {
		// tens of empty requests per minute, mah...
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var rc int
	var errS string
	var data string

	cacheKey := tgt.Service + "://" + tgt.HostIp
	if skip {
		cacheKey += "/skip"
	}

	var cacheOk bool
	if withCache {
		cReply, ok := hbcache.Check(cacheKey)
		if ok {
			rc = int(cReply.Status)
			errS = cReply.Error
			data = cReply.Data
			cacheOk = true
		}
	}

	if !withCache || !cacheOk {
		out, err := heartbleed.Heartbleed(tgt, PAYLOAD, skip)

		if err == heartbleed.Safe || err == heartbleed.Closed {
			rc = 1
		} else if err != nil {
			rc = 2
		} else {
			rc = 0
			// _, err := bleed.Heartbleed(tgt, PAYLOAD)
			// if err == nil {
			// 	// Two VULN in a row
			// 	rc = 0
			// } else {
			// 	// One VULN and one not
			// 	_, err := bleed.Heartbleed(tgt, PAYLOAD)
			// 	if err == nil {
			// 		// 2 VULN on 3 tries
			// 		rc = 0
			// 	} else {
			// 		// 1 VULN on 3 tries
			// 		if err == bleed.Safe {
			// 			rc = 1
			// 		} else {
			// 			rc = 2
			// 		}
			// 	}
			// }
		}

		switch rc {
		case 0:
			data = out
			log.Printf("%v (%v) - VULNERABLE [skip: %v]", tgt.HostIp, tgt.Service, skip)
		case 1:
			log.Printf("%v (%v) - SAFE", tgt.HostIp, tgt.Service)
		case 2:
			errS = err.Error()
			if errS == "Please try again" {
				log.Printf("%v (%v) - MISMATCH", tgt.HostIp, tgt.Service)
			} else {
				log.Printf("%v (%v) - ERROR [%v]", tgt.HostIp, tgt.Service, errS)
			}
		}
	}

	if withCache && !cacheOk {
		hbcache.Set(cacheKey, rc, data, errS)
	}

	res := result{rc, data, errS, tgt.HostIp}
	j, err := json.Marshal(res)
	if err != nil {
		log.Println("[json] ERROR:", err)
	} else {
		w.Write(j)
	}
}

func bleedHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Path[len("/bleed/"):]

	tgt := heartbleed.Target{
		HostIp:  string(host),
		Service: "https",
	}
	handleRequest(&tgt, w, r, true)
}

func bleedQueryHandler(w http.ResponseWriter, r *http.Request) {
	q, ok := r.URL.Query()["u"]
	if !ok || len(q) != 1 {
		return
	}

	skip, ok := r.URL.Query()["skip"]
	s := false
	if ok && len(skip) == 1 {
		s = true
	}

	tgt := heartbleed.Target{
		HostIp:  string(q[0]),
		Service: "https",
	}

	u, err := url.Parse(tgt.HostIp)
	if err == nil && u.Host != "" {
		tgt.HostIp = u.Host
		if u.Scheme != "" {
			tgt.Service = u.Scheme
		}
	}

	handleRequest(&tgt, w, r, s)
}

var Usage = `Heartbleed test server.

Usage:
  HBserver --redir-host=<host> [--listen=<addr:port> --expiry=<duration>]
           [--key=<key> --cert=<cert>]
  HBserver -h | --help
  HBserver --version

Options:
  --redir-host HOST   Redirect requests to "/" to this host.
  --listen ADDR:PORT  Listen and serve requests to this address:port [default: :8082].
  --expiry DURATION   ENABLE CACHING. Expire records after this period.
                      Uses Go's parse syntax
                      e.g. 10m = 10 minutes, 600s = 600 seconds, 1d = 1 day, etc.
  --key KEY           TLS key .pem file -- enable TLS
  --cert CERT         TLS cert .pem file -- enable TLS
  -h --help           Show this screen.
  --version           Show version.`

func main() {
	arguments, _ := docopt.Parse(Usage, nil, true, "HBserver 0.3", false)

	if arguments["--expiry"] != nil {
		withCache = true
	}

	if withCache {
		hbcache.Init(arguments["--expiry"].(string))
	}

	http.Handle("/", http.RedirectHandler(
		arguments["--redir-host"].(string), http.StatusFound,
	))

	// Required for some ELBs
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/bleed/", bleedHandler)
	http.HandleFunc("/bleed/query", bleedQueryHandler)

	s := &http.Server{
		Addr:        arguments["--listen"].(string),
		ReadTimeout: 10 * time.Second,
	}

	if arguments["--key"] != nil && arguments["--cert"] != nil {
		log.Printf("Starting TLS server on %s\n", s.Addr)
		log.Fatal("ListenAndServeTLS: ", s.ListenAndServeTLS(
			arguments["--cert"].(string), arguments["--key"].(string),
		))
	} else {
		log.Printf("Starting server on %s\n", s.Addr)
		log.Fatal("ListenAndServe: ", s.ListenAndServe())
	}
}
