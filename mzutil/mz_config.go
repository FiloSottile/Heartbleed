/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package util

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

/* Craptastic typeless parser to read config values (use until things
   settle and you can use something more efficent like TOML)
*/

type JsMap map[string]interface{}

type MzConfig struct {
	config JsMap
	flags  map[string]bool
}

/* Read a ini like configuration file into a map
 */
func ReadMzConfig(filename string) (config *MzConfig, err error) {
	// Yay for no equivalent to readln
	config = &MzConfig{
		config: make(JsMap),
		flags:  make(map[string]bool),
	}
	file, err := os.Open(filename)

	defer file.Close()

	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
		// skip lines beginning with '#/;'
		if strings.Contains("#/;", string(line[0])) {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) < 2 {
			continue
		}
		config.config[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	if err != nil && err != io.EOF {
		log.Panic(err)
	}
	return config, nil
}

/* Get a value from the config map, providing an optional default.
   This is a fairly common behavior for me.
*/
func (self *MzConfig) Get(key string, def string) string {
	if val, ok := self.config[key]; ok {
		return val.(string)
	}
	return def
}

/* Set a value if it's not already defined
 */
func (self *MzConfig) SetDefault(key string, val string) string {
	if _, ok := self.config[key]; !ok {
		self.config[key] = val
	}
	return self.config[key].(string)
}

/* Test for a boolean flag. Missing flags are false.
 */
func (self *MzConfig) GetFlag(key string) bool {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	if flag, ok := self.flags[key]; ok {
		return flag
	}
	if val, ok := self.config[key]; ok {
		self.flags[key], _ = strconv.ParseBool(val.(string))
		return self.flags[key]
	}
	return false
}

/* Set the boolean flag if not already specified
 */
func (self MzConfig) SetDefaultFlag(key string, val bool) (flag bool) {
	if bflag, ok := self.flags[key]; ok {
		return bflag
	}
	if _, ok := self.config[key]; ok {
		return self.GetFlag(key)
	}
	self.flags[key] = val
	return val
}

// o4fs
// vim: set tabstab=4 softtabstop=4 shiftwidth=4 noexpandtab
