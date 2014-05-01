package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/docopt/docopt-go"

	bleed "github.com/FiloSottile/Heartbleed/bleed"
	cache "github.com/FiloSottile/Heartbleed/server/cache"
)

var (
	PAYLOAD   = []byte("filippo.io/Heartbleed")
	withCache = false
	cacheData = false
)

type result struct {
	Code  int    `json:"code"`
	Data  string `json:"data"`
	Error string `json:"error"`
	Host  string `json:"host"`
}

func handleRequest(tgt *bleed.Target, w http.ResponseWriter, r *http.Request, skip bool) {
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
		cReply, ok := cache.Check(cacheKey)
		if ok {
			rc = int(cReply.Status)
			errS = cReply.Error
			data = cReply.Data
			cacheOk = true
		}
	}

	if !withCache || !cacheOk {
		out, err := bleed.Heartbleed(tgt, PAYLOAD, skip)

		if err == bleed.Safe || err == bleed.Closed {
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
		/* Storing the data returned from the site is problematic for a
		   number of reasons. While it may be valuable in certain
		   circumstances, it's not a good idea to store and return this
		   data to arbitrary parties. */
		cdata := ""
		if cacheData {
			cdata = data
		}
		cache.Set(cacheKey, rc, cdata, errS)
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

	tgt := bleed.Target{
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

	tgt := bleed.Target{
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
  HBserver --redir-host=<host> [--listen=<addr:port> --expiry=<duration> --cache-data]
  HBserver -h | --help
  HBserver --version

Options:
  --redir-host HOST   Redirect requests to "/" to this host.
  --listen ADDR:PORT  Listen and serve requests to this address:port [default: :8082].
  --expiry DURATION   ENABLE CACHING. Expire records after this period.
                      Uses Go's parse syntax
                      e.g. 10m = 10 minutes, 600s = 600 seconds, 1d = 1 day, etc.
  -h --help           Show this screen.
  --cache-data        Cache the data. (Not recommended. May contain private
                      info or other unwanted data.)
  --version           Show version.`

func main() {
	arguments, _ := docopt.Parse(Usage, nil, true, "HBserver 0.3", false)

	if arguments["--expiry"] != nil {
		withCache = true
	}

	if arguments["--cache-data"] != nil {
		cacheData = true
	}

	// this was returning nil in later code. Best be certain it's defined.
	if _, ok := arguments["--listen"]; !ok {
		arguments["--listen"] = ":8082"
	}

	if withCache {
		cache.Init(arguments["--expiry"].(string))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, arguments["--redir-host"].(string), http.StatusFound)
	})

	// Required for some ELBs
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/bleed/", bleedHandler)
	http.HandleFunc("/bleed/query", bleedQueryHandler)

	log.Printf("Starting server on %s\n", arguments["--listen"].(string))
	err := http.ListenAndServe(arguments["--listen"].(string), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
