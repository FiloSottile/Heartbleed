// Used internally by GoDynamo for nullable primitives
package nullable

import (
	"encoding/json"
)

// NullableString is a string that when empty is marshaled as null
// use this when there is an *optional* string parameter.
type NullableString string

func (n NullableString) MarshalJSON() ([]byte, error) {
	sn := string(n)
	if sn == "" {
		var i interface{}
		return json.Marshal(i)
	}
	return json.Marshal(sn)
}

// NullableUInt64 is a uint64 that when empty (0) is marshaled as null
// use this when there is an *optional* uint64 parameter
type NullableUInt64 uint64

func (n NullableUInt64) MarshalJSON() ([]byte, error) {
	in := uint64(n)
	if in == 0 {
		var i interface{}
		return json.Marshal(i)
	}
	return json.Marshal(in)
}
