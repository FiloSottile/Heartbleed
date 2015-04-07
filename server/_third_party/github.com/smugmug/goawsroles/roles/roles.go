// Defines an interface that can be implemented to provide IAM Roles data through various communication
// mechanisms, most likely regular text files (see the roles_files.go implementation).
package roles

type RolesFields struct {
	AccessKey string
	Secret    string
	Token     string
}

// NewRolesFields returns a pointer to a RolesFields instance.
func NewRolesFields() *RolesFields {
	return new(RolesFields)
}

// IsEmpty determines if a RolesField struct is uninitialized.
func (rf *RolesFields) IsEmpty() bool {
	return rf.AccessKey == "" || rf.Secret == "" || rf.Token == ""
}

// ZeroRoles recreate the RolessFields as initialized by NewRolesFields.
func (rf *RolesFields) ZeroRoles() {
	rf.AccessKey = ""
	rf.Secret = ""
	rf.Token = ""
}

// RolesReader is our interface to describe the functionality for roles credential information.
type RolesReader interface {
	// blocking read of roles from roles provider
	RolesRead() error
	// zero out roles values
	ZeroRoles()
	// test for emptiness
	IsEmpty() bool
	// getters
	// wrapper to individual getters
	Get() (string, string, string, error)
	// below funcs should be called by GetAllRoles
	GetAccessKey() (string, error)
	GetSecret() (string, error)
	GetToken() (string, error)
	// mechanism by which roles strings can be refreshed (event-based, polling etc).
	// if you don't want to implement this, just use it to wrap RolesRead
	RolesWatch(c chan error, s chan bool)
}
