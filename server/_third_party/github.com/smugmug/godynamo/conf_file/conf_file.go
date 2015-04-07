// Manages reading the conf file into the global var as described in the `conf` package.
package conf_file

import (
	"encoding/json"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/goawsroles/roles_files"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/aws_const"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/conf"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

// Read will look for and read in the conf file, which can then be referenced as conf.Vals.
// The conf file is specifically relevant to properly formatted requests, so it is currently
// called in the initialization of the authreq package.
// You may also set $GODYNAMO_CONF_FILE to be a fully-qualified path to a conf file
// if the two preset locations are not adequate.
func Read() {
	var cf conf.SDK_conf_file
	local_conf := os.Getenv("HOME") + string(filepath.Separator) + "." + conf.CONF_NAME
	etc_conf := string(filepath.Separator) + "etc" + string(filepath.Separator) + conf.CONF_NAME
	read_conf := false
	conf_files := make([]string, 0)
	// assumes that if set, this is a fully-qualified file path
	const env_conf = "GODYNAMO_CONF_FILE"
	// assumes that if set, this is a fully-qualified file path
	if os.Getenv(env_conf) != "" {
		conf_files = append(conf_files, os.Getenv(env_conf))
	}
	conf_files = append(conf_files, local_conf)
	conf_files = append(conf_files, etc_conf)
	conf.Vals.ConfLock.Lock()
	defer conf.Vals.ConfLock.Unlock()
CONF_LOCATIONS:
	for _, conf_file := range conf_files {
		conf_bytes, conf_err := ioutil.ReadFile(conf_file)
		if conf_err != nil {
			log.Printf("cannot find conf file at %s\n", conf_file)
			continue CONF_LOCATIONS
		} else {
			um_err := json.Unmarshal(conf_bytes, &cf)
			if um_err != nil {
				panic("conf_file.Read:" + conf_file +
					" json err: " +
					um_err.Error())
			} else {
				log.Printf("read conf from: %s\n", conf_file)
				read_conf = true
				break
			}
		}
	}
	if !read_conf {
		panic("confload.init: read err: " +
			"\n\n\n*****\nMake sure you have a conf file!\n" +
			"An example conf file is located in the /conf dir.\n" +
			"Put it in your home dir as\n$HOME/.aws-config.json\nor " +
			"in /etc as\n/etc/aws-config.json\nand fill " +
			"in the values for your AWS account*****\n\n\n")
	}

	// make sure the dynamo endpoint is available
	addrs, addrs_err := net.LookupIP(cf.Services.Dynamo_db.Host)
	if addrs_err != nil {
		panic("cannot look up hostname: " + cf.Services.Dynamo_db.Host)
	}
	dynamo_ip := (addrs[0]).String()

	// assign the values to our globally-available conf.Vals struct instance
	conf.Vals.Auth.AccessKey = cf.Services.Default_settings.Params.Access_key_id
	conf.Vals.Auth.Secret = cf.Services.Default_settings.Params.Secret_access_key
	conf.Vals.UseSysLog = cf.Services.Default_settings.Params.Use_sys_log
	conf.Vals.Network.DynamoDB.Host = cf.Services.Dynamo_db.Host
	conf.Vals.Network.DynamoDB.IP = dynamo_ip
	conf.Vals.Network.DynamoDB.Zone = cf.Services.Dynamo_db.Zone
	scheme := "http"
	port := aws_const.PORT // already a string
	if cf.Services.Dynamo_db.Scheme != "" {
		scheme = cf.Services.Dynamo_db.Scheme
	}
	if cf.Services.Dynamo_db.Port != 0 {
		port = strconv.Itoa(cf.Services.Dynamo_db.Port)
	}
	conf.Vals.Network.DynamoDB.Port = port
	conf.Vals.Network.DynamoDB.Scheme = scheme
	conf.Vals.Network.DynamoDB.URL = scheme + "://" + conf.Vals.Network.DynamoDB.Host +
		":" + port
	_, url_err := url.Parse(conf.Vals.Network.DynamoDB.URL)
	if url_err != nil {
		panic("confload.init: read err: conf.Vals.Network.DynamoDB.URL malformed")
	}

	// If set to true, programs that are written with godynamo may
	// opt to launch the keepalive goroutine to keep conns open.
	conf.Vals.Network.DynamoDB.KeepAlive = cf.Services.Dynamo_db.KeepAlive

	// read in flags for IAM support
	if cf.Services.Dynamo_db.IAM.Use_iam == true {
		if cf.Services.Dynamo_db.IAM.Role_provider != roles_files.ROLE_PROVIDER {
			panic("confload.init: read err: " +
				"\n\n\n**** only IAM role provider 'file' is supported *****\n\n\n")
		}
		conf.Vals.IAM.RoleProvider = cf.Services.Dynamo_db.IAM.Role_provider
		conf.Vals.IAM.File.BaseDir = cf.Services.Dynamo_db.IAM.Base_dir
		conf.Vals.IAM.File.AccessKey = cf.Services.Dynamo_db.IAM.Access_key
		conf.Vals.IAM.File.Secret = cf.Services.Dynamo_db.IAM.Secret_key
		conf.Vals.IAM.File.Token = cf.Services.Dynamo_db.IAM.Token
		if cf.Services.Dynamo_db.IAM.Watch == true {
			conf.Vals.IAM.Watch = true
		} else {
			conf.Vals.IAM.Watch = false
		}
		conf.Vals.UseIAM = true
	}
	conf.Vals.Initialized = true
}
