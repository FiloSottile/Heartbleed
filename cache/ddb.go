/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

// Dynamo DB cache layer

/* Depends on godynamo. Please refer to that package for how to configure
 * this package.
 *
 * The dynamo table we use has the following structure:
 *  primary Hash key: hostname
 *      No range key
 *      No secondary indexes.
 *  columns:
 *      Mtime   - last access time
 *      Status  - status return value
 *
 */
package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	ep "github.com/smugmug/godynamo/endpoint"
	get "github.com/smugmug/godynamo/endpoints/get_item"
	put "github.com/smugmug/godynamo/endpoints/put_item"
	keepalive "github.com/smugmug/godynamo/keepalive"
)

var CACHE_TAB = "mozHeartbleed"
var expry time.Duration

type CacheReply struct {
	Host       string
	LastUpdate int64
	Status     int64
}

/* Initialize the cache layer.
 * configPath - Path to the godynamo JSON config file.
 * expiration - Go standard duration string, defaults to '10m'
 *
 * This forces the Go Dynamo Config path into an environment var
 * so that its preferred config reader can find it. Leave blank if
 * you want this to use a config file in the default places.
 */
func Init(configPath string, expiration string) (err error) {
	os.Setenv("GODYNAMO_CONF_FILE", configPath)
	conf_file.Read()

	if conf.Vals.Initialized == false {
		panic("Uninitialized conf.Vals global")
	}

	if conf.Vals.Network.DynamoDB.KeepAlive {
        log.Printf("Cache Info: Launching background DynamoDB keepalive")
		go keepalive.KeepAlive([]string{conf.Vals.Network.DynamoDB.URL})
	}

	expry, err = time.ParseDuration(expiration)
    if err != nil {
        log.Printf("Cache Warn: Could not parse expry string. Expiring after 10m [%s]", err.Error())
        expry = time.Minute * 10
    }

	// if we were using IAM, put that code here.
	return
}

/* Fetch the record from the Cache.
 *
 * returns the record and if the record should be considered "OK" to use.
 * An OK record is a valid, non-expired cache stored entry.
 */
func Check(host string) (reply CacheReply, ok bool) {
	var getr get.Request
	var gr get.Response

	ok = true
	getr.TableName = CACHE_TAB
	getr.Key = make(ep.Item)
	getr.Key["hostname"] = ep.AttributeValue{S: host}
	body, code, err := getr.EndpointReq()
	if err != nil || code != http.StatusOK {
		if err != nil {
			log.Printf("Cache Error: %s\n", err.Error())
		}
		ok = false
		return
	}
	// get the time from the body,
	//log.Printf("####CACHE_GET: %s %d\n", string(body), len(body))
	if len(body) < 3 {
        // record not found
		ok = false
		return reply, ok
	}

	if err = json.Unmarshal([]byte(body), &gr); err == nil {
		reply.LastUpdate, err = strconv.ParseInt(gr.Item["Mtime"].N, 10, 64)
		if err != nil {
            // unparsable time
			log.Printf("Cache Error: Bad Record %s, %s", host, err.Error())
			ok = false
			return
		}
        if reply.LastUpdate < time.Now().UTC().Truncate(expry).Unix(){
            // record has expired.
            ok = false
        }
		reply.Status, err = strconv.ParseInt(gr.Item["Status"].N, 10, 64)
		if err != nil {
            //unparsable status
			log.Printf("Cache Error: Bad Record %s, %s", host, err.Error())
			ok = false
			return
		}
		reply.Host = gr.Item["hostname"].S
		//log.Printf("gr %+v", reply)
	}
	return
}

/* Set the state of the host in the cache
 *
 */
func Set(host string, state int) (err error) {
	var putr put.Request
	//var status string

	putr.TableName = CACHE_TAB
	putr.Item = make(ep.Item)
	putr.Item["hostname"] = ep.AttributeValue{S: host}
	putr.Item["Mtime"] = ep.AttributeValue{N: strconv.FormatInt(time.Now().UTC().Unix(), 10)}
	putr.Item["Status"] = ep.AttributeValue{N: strconv.FormatInt(int64(state), 10)}
	body, code, err := putr.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("Cache Error: put failed %d, %v, %s", code, err, body)
	}
	//log.Printf("####CACHE_SET: %s\n", string(body))
	return
}
