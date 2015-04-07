// Manages reading the configuration file format into our shared conf state.
// The AWS SDKs utilize a conf file format that this package attempts compatibility with while
// supporting extra fields. This is why the type detailing the file and internal formats differ.
// See SAMPLE-aws-config.json in the source repository for a sample.
package conf

import (
	"sync"
)

const (
	CONF_NAME          = "aws-config.json"
	ROLE_PROVIDER_FILE = "file"
)

// SDK_conf_File roughly matches the format as used by recent amazon SDKs, plus some additions.
// These are correlated to the fields you would fill in the conf file
type SDK_conf_file struct {
	Extends  []string
	Services struct {
		Default_settings struct {
			Params struct {
				// Traditional AWS access/secret authentication pair.
				Access_key_id     string
				Secret_access_key string
				// If you use syslogd (a linux or *bsd system), you may set this to "true".
				// (this is currently unused)
				Use_sys_log bool
			}
		}
		Dynamo_db struct {
			// Your dynamo hostname.
			Host string
			// Typically http or https, will have "://" appended.
			Scheme string
			// Port should correspond to the scheme.
			Port int
			// If set to true, programs that are written with godynamo may
			// opt to launch the keepalive goroutine to keep conns open.
			KeepAlive bool
			// Your aws zone.
			Zone string
			IAM  struct {
				// Set to true to use IAM authentication.
				Use_iam bool
				// The role provider is described in the goawsroles package.
				// See: https://github.com/smugmug/goawsroles/
				// Currently the only support is for the "file" provider, whereby
				// roles data is written to local files.
				Role_provider string
				// The identifier (filename, etc) for the IAM Access Key
				Access_key string
				// The identifier (filename, etc) for the IAM Secret Key
				Secret_key string
				// The identifier (filename, etc) for the IAM Token
				Token string
				// If using the "file" role provider, the base dir to read IAM files.
				Base_dir string
				// Set to true if you would like the roles resource watched for changes
				// and automatically (and atomically) updated.
				Watch bool
			}
		}
	}
}

// AWS_Conf is the structure used internally in godynamo.
type AWS_Conf struct {
	// Set to true if this struct is populated correctly.
	Initialized bool
	// Traditional AWS authentication pair.
	Auth struct {
		AccessKey string
		Secret    string
	}
	// Dynamo connection data.
	Network struct {
		DynamoDB struct {
			Host   string
			Scheme string
			// Port is converted into a string for internal use, typically
			// stitching together URL path strings.
			Port      string
			KeepAlive bool
			IP        string
			Zone      string
			URL       string
		}
	}
	// If using syslogd
	UseSysLog bool
	// If using IAM
	UseIAM bool
	// The IAM role provider info
	IAM struct {
		RoleProvider string
		Watch        bool
		// Tells you where the credentials can be read from
		File struct {
			AccessKey string
			Secret    string
			Token     string
			BaseDir   string
		}
		// The credentials themselves, once loaded from Files.* above
		// these are kept distinct from the global AccessKey and Secret
		// in the event a caller wants a mixed model
		Credentials struct {
			AccessKey string
			Secret    string
			Token     string
		}
	}
	// Lock used when accessing IAM values, which will change during execution.
	// other values will persist for program duration so they can be read without locking.
	ConfLock sync.RWMutex
}

// Vals is the global conf vals struct. It is shared throughout the duration of program execution.
// Use the embedded ConfLock mutex to use it safely.
var (
	Vals AWS_Conf
)
