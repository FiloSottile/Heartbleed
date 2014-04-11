//+build skip

// From heartbleed.fillip.io
// Adding DynamoDB caching.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	bleed "github.com/FiloSottile/Heartbleed/bleed"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	ep "github.com/smugmug/godynamo/endpoint"
	get "github.com/smugmug/godynamo/endpoints/get_item"
	put "github.com/smugmug/godynamo/endpoints/put_item"
	keepalive "github.com/smugmug/godynamo/keepalive"
)

var PAYLOAD = []byte("heartbleed.mozilla.com")
var LOCALHOST = "http://localhost"
var PORT_SRV = ":8082"
var CACHE_TAB = "mozHeartbleed"
var EXPRY time.Duration

const (
	VUNERABLE = iota
	SAFE
	ERROR
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, LOCALHOST, http.StatusFound)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

type result struct {
	Code  int    `json:"code"`
	Data  string `json:"data"`
	Error string `json:"error"`
}

type cacheReply struct {
	Host       string
	LastUpdate int64
	Status     int64
}

type jrep map[string]interface{}

func cacheCheck(host string) (reply cacheReply, ok bool) {
	var getr get.Request
	var gr get.Response

	ok = false
	getr.TableName = CACHE_TAB
	getr.Key = make(ep.Item)
	getr.Key["hostname"] = ep.AttributeValue{S: host}
	body, code, err := getr.EndpointReq()
	if err != nil || code != http.StatusOK {
		if err != nil {
			log.Printf("!!!! Error: %s\n", err.Error())
		}
		ok = false
		return
	}
	// get the time from the body,
	//log.Printf("####CACHE_GET: %s %d\n", string(body), len(body))
	if len(body) < 3 {
		ok = false
	}

	if err = json.Unmarshal([]byte(body), &gr); err == nil {
		reply.LastUpdate, err = strconv.ParseInt(gr.Item["Mtime"].N, 10, 64)
		if err != nil {
			log.Printf("Bad Record %s", host)
			ok = false
			return
		}
		reply.Status, err = strconv.ParseInt(gr.Item["Status"].N, 10, 64)
		if err != nil {
			log.Printf("Bad Record %s", host)
			ok = false
			return
		}
		reply.Host = gr.Item["hostname"].S
		//log.Printf("gr %+v", reply)
		ok = true
	}
	// if the time has not expired, then things are good.
	// else retry.
	return
}

func cacheSet(host string, state int) (err error) {
	var putr put.Request
	//var status string

	putr.TableName = CACHE_TAB
	putr.Item = make(ep.Item)
	putr.Item["hostname"] = ep.AttributeValue{S: host}
	putr.Item["Mtime"] = ep.AttributeValue{N: strconv.FormatInt(time.Now().UTC().Unix(), 10)}
	putr.Item["Status"] = ep.AttributeValue{N: strconv.FormatInt(int64(state), 10)}
	body, code, err := putr.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("put failed %d, %v, %s", code, err, body)
	}
	//log.Printf("####CACHE_SET: %s\n", string(body))
	return
}

func bleedHandler(w http.ResponseWriter, r *http.Request) {
	var rc int
	var ok bool
	var fullCheck bool
	var cReply cacheReply
	var data []byte
	var err error
	var errS string

	w.Header().Set("Access-Control-Allow-Origin", "*")
	host := r.URL.Path[len("/bleed/"):]
	u, err := url.Parse(host)
	if err == nil && u.Host != "" {
		host = u.Host
	}

	tgt := bleed.Target{
		HostIp:  string(host),
		Service: "https",
	}

	// check cache
	fullCheck = true
	data = []byte("")
	if cReply, ok = cacheCheck(tgt.HostIp); ok {
		if cReply.LastUpdate < time.Now().UTC().Truncate(EXPRY).Unix() {
			log.Printf("Refetching " + tgt.HostIp)
			rc = int(cReply.Status)
		} else {
			fullCheck = false
		}
	}
	if fullCheck {
		data, err = bleed.Heartbleed(&tgt, PAYLOAD)
		if err == bleed.Safe {
			rc = SAFE
			data = []byte("")
			log.Printf("%v - SAFE", host)
		} else if err != nil {
			rc = ERROR
			data = []byte("")
			errS = err.Error()
			log.Printf("%v - ERROR", host)
		} else {
			rc = VUNERABLE
			log.Printf("%v - VULNERABLE", host)
		}
		// record cache
		err = cacheSet(host, rc)
	}
	res := result{rc, string(data), errS}
	j, err := json.Marshal(res)
	if err != nil {
		log.Println("ERROR", err)
	} else {
		w.Write(j)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {

	var err error

	// should take a conf arg
	conf_file.Read()
	EXPRY, err = time.ParseDuration("10m")

	if conf.Vals.Initialized == false {
		panic("Uninitialized conf.Vals global")
	}

	if conf.Vals.Network.DynamoDB.KeepAlive {
		log.Printf("Launching background DynamoDB keepalive")
		go keepalive.KeepAlive([]string{conf.Vals.Network.DynamoDB.URL})
	}
	// if we were using IAM, put that code here.

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/bleed/", bleedHandler)
	log.Printf("Starting server on %s\n", PORT_SRV)
	err = http.ListenAndServe(PORT_SRV, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
