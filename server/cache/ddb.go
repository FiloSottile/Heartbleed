/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

// Dynamo DB cache layer

/* Depends on godynamo. Please refer to that package for how to configure
 * this package.
 *
 * The dynamo table we use has the following structure:
 *  Primary Hash key: hostname
 *  No range key
 *  No secondary indexes.
 *  columns:
 *      Mtime   - last access time
 *      Status  - status return value
 *      Data    - memory dump string
 *      Error   - error string value
 *
 */
package hbcache

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/conf"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/conf_file"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/endpoints/get_item"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/endpoints/put_item"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/keepalive"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributevalue"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/item"
)

var CACHE_TAB = "Heartbleed"
var expiry time.Duration

type CacheReply struct {
	Host       string
	LastUpdate int64
	Status     int64
	Data       string
	Error      string
}

/* Initialize the cache layer.
 * expiration - Go standard duration string, defaults to '10m'
 */
func Init(expiration string) {
	conf_file.Read()

	if conf.Vals.Initialized == false {
		panic("Uninitialized conf.Vals global")
	}

	if conf.Vals.Network.DynamoDB.KeepAlive {
		log.Printf("[cache] INFO: Launching background DynamoDB keepalive")
		go keepalive.KeepAlive([]string{conf.Vals.Network.DynamoDB.URL})
	}

	exp, err := time.ParseDuration(expiration)
	if err != nil {
		log.Printf("[cache] WARN: Could not parse expiry string [%s]. Expiring after 10m.", err.Error())
		expiry = time.Minute * 10
	} else {
		expiry = exp
	}
}

/* Fetch the record from the Cache.
 *
 * returns the record and if the record should be considered "OK" to use.
 * An OK record is a valid, non-expired cache stored entry.
 */
func Check(host string) (CacheReply, bool) {
	getr := get_item.Request{
		TableName: CACHE_TAB,
		Key:       make(item.Key),
	}
	getr.Key["hostname"] = &attributevalue.AttributeValue{S: host}
	body, code, err := getr.EndpointReq()
	if err != nil {
		log.Printf("[cache] ERROR: %s", err.Error())
		return CacheReply{}, false
	}
	if code != http.StatusOK {
		return CacheReply{}, false
	}

	// get the time from the body
	// log.Printf("####CACHE_GET: %s %d\n", string(body), len(body))
	if len(body) < 3 {
		// record not found
		return CacheReply{}, false
	}

	var gr get_item.Response
	var reply CacheReply
	if err = json.Unmarshal([]byte(body), &gr); err == nil {
		reply.Status, err = strconv.ParseInt(gr.Item["Status"].N, 10, 64)
		if err != nil {
			// unparsable status
			log.Printf("[cache] ERROR: Bad Record %s, %s", host, err.Error())
			return CacheReply{}, false
		}

		reply.Host = gr.Item["hostname"].S
		reply.Error = gr.Item["Error"].S
		if reply.Error == "---" {
			reply.Error = ""
		}
		reply.Data = gr.Item["Data"].S
		if reply.Data == "---" {
			reply.Data = ""
		}

		reply.LastUpdate, err = strconv.ParseInt(gr.Item["Mtime"].N, 10, 64)
		if err != nil {
			// unparsable time
			log.Printf("[cache] ERROR: Bad Record %s, %s", host, err.Error())
			return CacheReply{}, false
		}
		if reply.LastUpdate < time.Now().UTC().Add(-expiry).Unix() {
			// record has expired.
			return reply, false
		}
		// log.Printf("gr %+v", reply)
	}

	return reply, true
}

/* Set the state of the host in the cache
 *
 */
func Set(host string, state int, data, errS string) error {
	if data == "" {
		data = "---"
	}
	if errS == "" {
		errS = "---"
	}

	putr := put_item.Request{
		TableName: CACHE_TAB,
		Item:      make(item.Item),
	}

	// Randomize the exp time by inserting the value up to "expiry" in the past
	// THIS MEANS THAT ITEMS ARE CACHED FOR A TIME FROM 0 TO expiry
	rnd := rand.Int63n(int64(expiry))
	mtime := time.Now().UTC().Add(-time.Duration(rnd)).Unix()

	putr.Item["hostname"] = &attributevalue.AttributeValue{S: host}
	putr.Item["Mtime"] = &attributevalue.AttributeValue{N: strconv.FormatInt(mtime, 10)}
	putr.Item["Status"] = &attributevalue.AttributeValue{N: strconv.FormatInt(int64(state), 10)}
	putr.Item["Data"] = &attributevalue.AttributeValue{S: data}
	putr.Item["Error"] = &attributevalue.AttributeValue{S: errS}

	body, code, err := putr.EndpointReq()
	if err != nil {
		log.Printf("[cache] ERROR: put failed %v, %s", err, body)
		return err
	}
	if code != http.StatusOK {
		log.Printf("[cache] ERROR: put failed %v, %s", code, body)
		return fmt.Errorf("put failed: code %d", code)
	}

	// log.Printf("####CACHE_SET: %s\n", string(body))
	return nil
}
