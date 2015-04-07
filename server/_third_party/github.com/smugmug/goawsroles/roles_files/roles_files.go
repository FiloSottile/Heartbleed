// Implements the RolesReader interface (roles.go) for text files. It is assumed that users will have
// another process that obtains roles data from AWS and deposits files.
//
// example use: see roles_files_test.go
//
package roles_files

import (
	"errors"
	"fmt"
	fsnotify "github.com/FiloSottile/Heartbleed/server/_third_party/github.com/howeyc/fsnotify"
	roles "github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/goawsroles/roles"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	FILE          = "file"
	ROLE_PROVIDER = FILE
)

// RolesFiles describes the location of roles files as well as a lock for safe access.
type RolesFiles struct {
	BaseDir       string
	AccessKeyFile string
	SecretFile    string
	TokenFile     string
	roleFields    *roles.RolesFields
	lock          sync.RWMutex
}

// NewRolesFiles returns a pointer to a RolesFields instance.
func NewRolesFiles() *RolesFiles {
	r := new(RolesFiles)
	r.roleFields = roles.NewRolesFields()
	return r
}

// IsEmpty determines if a RolesFiles struct is uninitialized.
func (rf *RolesFiles) IsEmpty() bool {
	return rf.AccessKeyFile == "" ||
		rf.SecretFile == "" ||
		rf.TokenFile == "" ||
		rf.roleFields.IsEmpty()
}

// ZeroRoles recreate the RolesFiles as initialized by NewRolesFiles
func (rf *RolesFiles) ZeroRoles() {
	rf.lock.Lock()
	defer rf.lock.Unlock()
	rf.AccessKeyFile = ""
	rf.SecretFile = ""
	rf.TokenFile = ""
	rf.roleFields.ZeroRoles()
}

// RolesRead populates rolesFields with blocking refresh of files
func (rf *RolesFiles) RolesRead() error {
	roles_err := rf.rolesFilesRead()
	if roles_err != nil {
		rf.ZeroRoles()
		return roles_err
	}
	return nil
}

// RolesWatch catches filesystem notify events to determine when new roles files are ready to be read
// in and used as new authentication values.
func (rf *RolesFiles) RolesWatch(err_chan chan error, read_signal chan bool) {
	log.Printf("initiate roles watching\n")
	watcher, watcher_err := fsnotify.NewWatcher()
	if watcher_err != nil {
		err_chan <- watcher_err
	}
	watch_err := watcher.Watch(rf.BaseDir)
	if watch_err != nil {
		err_chan <- watch_err
	}
	defer watcher.Close()
	touched_access_file := false
	touched_secret_file := false
	touched_token_file := false
	func() {
		for {
			select {
			case ev := <-watcher.Event:
				func() {
					if ev.IsModify() || ev.IsCreate() {
						ev_s := ev.String()
						// collect events for all of the role files.
						// we only want to read in and reset the
						// strings when all have been written - to
						// do so earlier would leave the strings
						// in an inconsistent state
						if strings.Contains(ev_s, rf.AccessKeyFile) {
							touched_access_file = true
						}
						if strings.Contains(ev_s, rf.SecretFile) {
							touched_secret_file = true
						}
						if strings.Contains(ev_s, rf.TokenFile) {
							touched_token_file = true
						}
						// once we have seen all of the role files trigger
						// events, we want to lock access, read them in,
						// unlock access, and unset our flags
						if touched_access_file &&
							touched_secret_file &&
							touched_token_file {
							// our existing perms should be adequate while
							// we provide for writes of new perms files to
							// finish.
							log.Printf("roles_files.RolesWatch: " +
								"sleep (1) to allow file ops to complete")
							time.Sleep(time.Duration(1) * time.Second)
							roles_err := rf.rolesFilesRead()
							if roles_err != nil {
								e := fmt.Sprintf("roles_files.RolesWatch: "+
									"zeroing all roles on err:%s",
									roles_err.Error())
								log.Printf(e)
								rf.ZeroRoles()
								err_chan <- roles_err
							} else {
								log.Printf("roles_files.RolesWatch: "+
									"succesful re-read on %s\n",
									ev_s)
								touched_access_file = false
								touched_secret_file = false
								touched_token_file = false
								read_signal <- true
							}
						}
					}
				}()
			case err := <-watcher.Error:
				// return an error to the caller via the channel
				err_chan <- err
			}
		}
	}()
	// typically not reached, but signifies normal termination
	log.Printf("terminating roles watching\n")
	err_chan <- nil
}

func (rf *RolesFiles) Get() (string, string, string, error) {
	accessKey, accessKey_err := rf.GetAccessKey()
	if accessKey_err != nil {
		return "", "", "", accessKey_err
	}
	secret, secret_err := rf.GetSecret()
	if accessKey_err != nil {
		return "", "", "", secret_err
	}
	token, token_err := rf.GetToken()
	if token_err != nil {
		return "", "", "", token_err
	}
	return accessKey, secret, token, nil
}

func (rf *RolesFiles) GetAccessKey() (string, error) {
	rf.lock.RLock()
	defer rf.lock.RUnlock()
	if rf.roleFields.AccessKey == "" {
		return "", errors.New("roles_files.GetAccessKey: empty AccessKey")
	} else {
		return rf.roleFields.AccessKey, nil
	}
}

func (rf *RolesFiles) GetSecret() (string, error) {
	rf.lock.RLock()
	defer rf.lock.RUnlock()
	if rf.roleFields.Secret == "" {
		return "", errors.New("roles_files.GetSecret: empty Secret")
	} else {
		return rf.roleFields.Secret, nil
	}
}

func (rf *RolesFiles) GetToken() (string, error) {
	rf.lock.RLock()
	defer rf.lock.RUnlock()
	if rf.roleFields.Token == "" {
		return "", errors.New("roles_files.GetToken: empty Token")
	} else {
		return rf.roleFields.Token, nil
	}
}

func role_file_bytes(role_file_path string) ([]byte, error) {
	role_file_bytes, role_file_err := ioutil.ReadFile(role_file_path)
	if role_file_err != nil || len(role_file_bytes) == 0 {
		fe := ""
		if role_file_err != nil {
			fe = role_file_err.Error()
		} else {
			fe = "empty file, no err msg"
		}
		e := fmt.Sprintf("roles_files.role_file_bytes: %s read err: %s",
			role_file_path, fe)
		if role_file_err != nil {
			e += " " + role_file_err.Error()
		}
		return nil, errors.New(e)
	}
	return role_file_bytes, nil
}

// safety check - the mtimes of the role files should be within sixty seconds of another
func valid_mtime_range(ts []time.Time) bool {
	uts := make([]int, len(ts))
	for i, _ := range ts {
		uts[i] = int(ts[i].Unix())
	}
	sort.Sort(sort.Reverse(sort.IntSlice(uts)))
	return ((uts[0] - uts[len(uts)-1]) < 10)
}

// will read in data from three role files and return the struct defined in the the interface
func (rf *RolesFiles) rolesFilesRead() error {
	rf.lock.Lock()
	defer rf.lock.Unlock()
	if rf.BaseDir == "" {
		e := fmt.Sprintf("roles_files.rolesFilesRead: must specify a non-empty BaseDir")
		return errors.New(e)
	}
	accessKey_path := rf.BaseDir + string(os.PathSeparator) + rf.AccessKeyFile
	accessKey_bytes, accessKey_err := role_file_bytes(accessKey_path)
	if accessKey_err != nil {
		return accessKey_err
	}
	secret_path := rf.BaseDir + string(os.PathSeparator) + rf.SecretFile
	secret_bytes, secret_err := role_file_bytes(secret_path)
	if secret_err != nil {
		return secret_err
	}
	token_path := rf.BaseDir + string(os.PathSeparator) + rf.TokenFile
	token_bytes, token_err := role_file_bytes(token_path)
	if token_err != nil {
		return token_err
	}

	// get mod times for all of the files
	uts := make([]time.Time, 0)
	for _, role_file_path := range []string{accessKey_path, secret_path, token_path} {
		role_file_stat, role_file_stat_err := os.Stat(role_file_path)
		if role_file_stat_err != nil {
			e := fmt.Sprintf("roles_files.rolesFilesRead: role_file stat err: %s",
				role_file_stat_err.Error())
			return errors.New(e)
		}
		uts = append(uts, role_file_stat.ModTime())
	}

	if valid_mtime_range(uts) {
		rf.roleFields.AccessKey = string(accessKey_bytes)
		rf.roleFields.Secret = string(secret_bytes)
		rf.roleFields.Token = string(token_bytes)
		log.Printf("roles_files.rolesFilesRead: succesful assignment of role data\n")
		return nil
	} else {
		e := fmt.Sprintf("roles_files.rolesFilesRead: range of mtimes of roles files >10 seconds")
		return errors.New(e)
	}
}
