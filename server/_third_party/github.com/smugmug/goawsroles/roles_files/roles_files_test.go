package roles_files

import (
	"fmt"
	roles "github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/goawsroles/roles"
	"testing"
)

func TestRolesFiles(t *testing.T) {
	rf_ := NewRolesFiles()
	if !rf_.IsEmpty() {
		t.Errorf("new RolesFile is not empty?")
	}
	rf_.BaseDir = "/etc/tags"
	rf_.AccessKeyFile = "role_access_key"
	rf_.SecretFile = "role_secret_key"
	rf_.TokenFile = "role_token"
	rf := roles.RolesReader(rf_)
	rr_err := rf.RolesRead()
	if rr_err != nil {
		t.Errorf(rr_err.Error())
	}
	access_key, access_key_err := rf.GetAccessKey()
	if access_key_err != nil {
		t.Errorf("cannot read access key?")
	} else {
		fmt.Printf("access key: %s\n", access_key)
	}
	secret, secret_err := rf.GetSecret()
	if secret_err != nil {
		t.Errorf("cannot read access key?")
	} else {
		fmt.Printf("secret: %s\n", secret)
	}
	token, token_err := rf.GetToken()
	if token_err != nil {
		t.Errorf("cannot read access key?")
	} else {
		fmt.Printf("token: %s\n", token)
	}

	// now unset the token file, which should cause all of our entries to be zero'd
	rf_.TokenFile = ""
	rf = roles.RolesReader(rf_)
	rr_err = rf.RolesRead()
	if rr_err == nil {
		t.Errorf("no err on unset token file?")
	}
	rf_.BaseDir = "/etc/tags"
	rf_.AccessKeyFile = "role_access_key"
	rf_.SecretFile = "role_secret_key"
	rf_.TokenFile = "role_token"
	rf = roles.RolesReader(rf_)

	// this will cause your program to appear to hang, you can comment the rest out to get a
	// completed test
	c := make(chan error)
	s := make(chan bool)
	go rf.RolesWatch(c, s)

	go func() {
		for {
			select {
			case <-s:
				fmt.Printf("*********** re-read the files\n")
				accessKey, secret, token, get_err := rf.Get()
				if get_err != nil {
					fmt.Printf("cannot get a role file:%s\n", get_err.Error())
				} else {
					fmt.Printf("access key:%s\nsecret:%s\ntoken:%s\n",
						accessKey, secret, token)
				}
			}
		}
	}()

	watch_err := <-c
	if watch_err != nil {
		fmt.Printf("error from watcher: %s\n", watch_err.Error())
	}

}
