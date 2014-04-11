//+build skip

package main

import (
	bleed "github.com/FiloSottile/Heartbleed/bleed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var PAYLOAD = []byte("heartbleed.filippo.io")

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://filippo.io/Heartbleed", http.StatusFound)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hi there.")
}

type result struct {
	Code  int    `json:"code"`
	Data  string `json:"data"`
	Error string `json:"error"`
	Host  string `json:"host"`
}

func handleRequest(tgt *bleed.Target, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, err := bleed.Heartbleed(tgt, PAYLOAD)
	var rc int
	var errS string
	if err == bleed.Safe {
		rc = 1
		data = []byte("")
		log.Printf("%v (%v) - SAFE", tgt.HostIp, tgt.Service)
	} else if err != nil {
		rc = 2
		data = []byte("")
		errS = err.Error()
		log.Printf("%v (%v) - ERROR", tgt.HostIp, tgt.Service)
	} else {
		rc = 0
		log.Printf("%v (%v) - VULNERABLE", tgt.HostIp, tgt.Service)
	}
	res := result{rc, string(data), errS, tgt.HostIp}
	j, err := json.Marshal(res)
	if err != nil {
		log.Println("ERROR", err)
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
	handleRequest(&tgt, w, r)
}

func bleedQueryHandler(w http.ResponseWriter, r *http.Request) {
	q, ok := r.URL.Query()["u"]
	if !ok || len(q) != 1 {
		return
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

	handleRequest(&tgt, w, r)
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/bleed/", bleedHandler)
	http.HandleFunc("/bleed/query", bleedQueryHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
