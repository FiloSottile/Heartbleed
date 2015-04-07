// Support for AttributeValue type. See
// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_AttributeValue.html
package attributevalue

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/cast"
	"strconv"
)

// SetList represents the SS,BS and NS types which are ostensibly sets but encoded as
// json lists. Duplicates are allowed but removed when marshaling or unmarshaling.
type SetList []string

// MarshalJSON will remove duplicates
func (s SetList) MarshalJSON() ([]byte, error) {
	m := make(map[string]bool)
	for _, v := range s {
		m[v] = true
	}
	t := make([]string, len(m))
	i := 0
	for k, _ := range m {
		t[i] = k
		i++
	}
	return json.Marshal(t)
}

// UnmarshalJSON will remove duplicates
func (s *SetList) UnmarshalJSON(data []byte) error {
	if s == nil {
		return errors.New("pointer receiver for unmarshal nil")
	}
	t := make([]string, 0)
	um_err := json.Unmarshal(data, &t)
	if um_err != nil {
		return um_err
	}
	m := make(map[string]bool)
	for _, v := range t {
		m[v] = true
	}
	for k, _ := range m {
		*s = append(*s, k)
	}
	return nil
}

type AttributeValue struct {
	N string `json:",omitempty"`
	S string `json:",omitempty"`
	B string `json:",omitempty"`

	// These are pointers so we can have a vacuous type (nil), otherwise, we don't know
	// if false was a set value or the default. To set these
	BOOL *bool `json:",omitempty"`
	NULL *bool `json:",omitempty"`

	L []*AttributeValue          `json:",omitempty"`
	M map[string]*AttributeValue `json:",omitempty"`

	SS SetList `json:",omitempty"`
	NS SetList `json:",omitempty"`
	BS SetList `json:",omitempty"`
}

// Empty determines if an AttributeValue is vacuous. Explicitly do not bother
// testing the boolean fields.
func (a *AttributeValue) Empty() bool {
	return a.N == "" &&
		a.S == "" &&
		a.B == "" &&
		len(a.M) == 0 &&
		len(a.L) == 0 &&
		len(a.SS) == 0 &&
		len(a.NS) == 0 &&
		len(a.BS) == 0 &&
		a.BOOL == nil &&
		a.NULL == nil
}

type attributevalue AttributeValue

// MarshalJSON will emit null if the AttributeValue is Empty
func (a AttributeValue) MarshalJSON() ([]byte, error) {
	if a.Empty() || (nil == &a) {
		return json.Marshal(nil)
	} else {
		return json.Marshal(attributevalue(a))
	}
}

// Valid determines if more than one field has been set (in which case it is invalid).
func (a *AttributeValue) Valid() bool {
	c := 0
	if a.S != "" {
		c++
		if c > 1 {
			return false
		}
	}
	if a.N != "" {
		c++
		if c > 1 {
			return false
		}
	}
	if a.B != "" {
		c++
		if c > 1 {
			return false
		}
	}
	if len(a.M) != 0 {
		c++
		if c > 1 {
			return false
		}
	}
	if len(a.L) != 0 {
		c++
		if c > 1 {
			return false
		}
	}
	if len(a.SS) != 0 {
		c++
		if c > 1 {
			return false
		}
	}
	if len(a.NS) != 0 {
		c++
		if c > 1 {
			return false
		}
	}
	if len(a.BS) != 0 {
		c++
		if c > 1 {
			return false
		}
	}
	if a.BOOL != nil {
		c++
		if c > 1 {
			return false
		}
	}
	if a.NULL != nil {
		c++
		if c > 1 {
			return false
		}
	}
	return true
}

func NewAttributeValue() *AttributeValue {
	a := new(AttributeValue)
	a.L = make([]*AttributeValue, 0)
	a.M = make(map[string]*AttributeValue)
	a.SS = make([]string, 0)
	a.NS = make([]string, 0)
	a.BS = make([]string, 0)

	// BOOL and NULL let to nil to represent vacuous state

	return a
}

// Copy makes a copy of the this AttributeValue.
func (a *AttributeValue) Copy(ac *AttributeValue) error {
	if ac == nil {
		return errors.New("copy target attributevalue instance is nil")
	}
	ac.S = a.S
	ac.N = a.N
	ac.B = a.B

	if a.BOOL == nil {
		ac.BOOL = nil
	} else {
		ac.BOOL = new(bool)
		*ac.BOOL = *a.BOOL
	}

	if a.NULL == nil {
		ac.NULL = nil
	} else {
		ac.NULL = new(bool)
		*ac.NULL = *a.NULL
	}

	l_ss := len(a.SS)
	if l_ss != 0 {
		ac.SS = make([]string, l_ss)
		copy(ac.SS, a.SS)
	}

	l_ns := len(a.NS)
	if l_ns != 0 {
		ac.NS = make([]string, l_ns)
		copy(ac.NS, a.NS)
	}

	l_bs := len(a.BS)
	if l_bs != 0 {
		ac.BS = make([]string, l_bs)
		copy(ac.BS, a.BS)
	}

	// L is a recursive type, so the copy must be recursive
	l_L := len(a.L)
	if l_L != 0 {
		ac.L = make([]*AttributeValue, l_L)
		for i, _ := range a.L {
			ac.L[i] = NewAttributeValue()
			L_i_cp_err := a.L[i].Copy(ac.L[i])
			if L_i_cp_err != nil {
				return L_i_cp_err
			}
		}
	}

	// M is a recursive type, so the copy must be recursive
	l_M := len(a.M)
	if l_M != 0 {
		ac.M = make(map[string]*AttributeValue, l_M)
		for k, _ := range a.M {
			ac.M[k] = NewAttributeValue()
			M_k_cp_err := a.M[k].Copy(ac.M[k])
			if M_k_cp_err != nil {
				return M_k_cp_err
			}
		}
	}
	return nil
}

// InsertS sets the S field to string k
func (a *AttributeValue) InsertS(k string) error {
	a.S = k
	return nil
}

// InsertN sets the N field to number string k
func (a *AttributeValue) InsertN(k string) error {
	fs, ferr := cast.AWSParseFloat(k)
	if ferr != nil {
		return ferr
	}
	a.N = fs
	return nil
}

// InsertN_float64 works like InsertN but takes a float64
func (a *AttributeValue) InsertN_float64(f float64) error {
	a.N = strconv.FormatFloat(f, 'f', -1, 64)
	return nil
}

// InsertB sets the B field to string k, which it is assumed the caller has
// already encoded.
func (a *AttributeValue) InsertB(k string) error {
	berr := cast.AWSParseBinary(k)
	if berr != nil {
		return berr
	}
	a.B = k
	return nil
}

// InsertB_unencoded adds a new plain string to the B field.
// The argument is assumed to be plaintext and will be base64 encoded.
func (a *AttributeValue) InsertB_unencoded(k string) error {
	a.B = base64.StdEncoding.EncodeToString([]byte(k))
	return nil
}

// InsertSS adds a new string to the ss (JSON: SS) set.
// SS is *generated* from an internal representation (UM_ss)
// as it transforms a map into a list (a "set")
func (a *AttributeValue) InsertSS(k string) error {
	for _, v := range a.SS {
		if v == k {
			return nil
		}
	}
	a.SS = append(a.SS, k)
	return nil
}

// InsertNS adds a new number string to the ns (JSON: NS) set.
// String is parsed to make sure it is a represents a valid float.
// NS is *generated* from an internal representation (UM_ns)
// as it transforms a map into a list (a "set")
func (a *AttributeValue) InsertNS(k string) error {
	fs, ferr := cast.AWSParseFloat(k)
	if ferr != nil {
		return ferr
	}
	for _, v := range a.NS {
		if v == fs {
			return nil
		}
	}
	a.NS = append(a.NS, fs)
	return nil
}

// InsertNS_float64 works like InsertNS but takes a float64
func (a *AttributeValue) InsertNS_float64(f float64) error {
	k := strconv.FormatFloat(f, 'f', -1, 64)
	for _, v := range a.NS {
		if v == k {
			return nil
		}
	}
	a.NS = append(a.NS, k)
	return nil
}

// InsertBS adds a new base64 string to the bs (JSON: BS) set.
// String is parsed to make sure it is a represents a valid base64 blob.
// BS is *generated* from an internal representation (UM_bs)
// as it transforms a map into a list (a "set").
// The argument is assumed to be already encoded by the caller.
func (a *AttributeValue) InsertBS(k string) error {
	berr := cast.AWSParseBinary(k)
	if berr != nil {
		return berr
	}
	for _, v := range a.BS {
		if v == k {
			return nil
		}
	}
	a.BS = append(a.BS, k)
	return nil
}

// InsertBS_unencoded adds a new plain string to the bs (JSON: BS) set.
// BS is *generated* from an internal representation (UM_bs)
// as it transforms a map into a list (a "set").
// The argument is assumed to be plaintext and will be base64 encoded.
func (a *AttributeValue) InsertBS_unencoded(k string) error {
	b64_k := base64.StdEncoding.EncodeToString([]byte(k))
	for _, v := range a.BS {
		if v == b64_k {
			return nil
		}
	}
	a.BS = append(a.BS, b64_k)
	return nil
}

// InsertL will append a pointer to a new AttributeValue v to the L list.
func (a *AttributeValue) InsertL(v *AttributeValue) error {
	v_cp := NewAttributeValue()
	cp_err := v.Copy(v_cp)
	if cp_err != nil {
		return cp_err
	}
	a.L = append(a.L, v_cp)
	return nil
}

// InsertM will insert a pointer to a new AttributeValue v to the M map for key k.
// If k was previously set in the M map, the value will be overwritten.
func (a *AttributeValue) InsertM(k string, v *AttributeValue) error {
	v_cp := NewAttributeValue()
	cp_err := v.Copy(v_cp)
	if cp_err != nil {
		return cp_err
	}
	a.M[k] = v_cp
	return nil
}

// InsertBOOL will set the BOOL field.
func (a *AttributeValue) InsertBOOL(b bool) error {
	if a.BOOL == nil {
		a.BOOL = new(bool)
	}
	*a.BOOL = b
	return nil
}

// InsertNULL will set the NULL field.
func (a *AttributeValue) InsertNULL(b bool) error {
	if a.NULL == nil {
		a.NULL = new(bool)
	}
	*a.NULL = b
	return nil
}

// AttributeValueMap is used throughout GoDynamo
type AttributeValueMap map[string]*AttributeValue

func NewAttributeValueMap() AttributeValueMap {
	m := make(map[string]*AttributeValue)
	return m
}

// AttributeValueUpdate is used in UpdateItem
type AttributeValueUpdate struct {
	Action string          `json:",omitempty"`
	Value  *AttributeValue `json:",omitempty"`
}

func NewAttributeValueUpdate() *AttributeValueUpdate {
	a := new(AttributeValueUpdate)
	a.Value = NewAttributeValue()
	return a
}

type AttributeValueUpdateMap map[string]*AttributeValueUpdate

func NewAttributeValueUpdateMap() AttributeValueUpdateMap {
	m := make(map[string]*AttributeValueUpdate)
	return m
}

// BasicJSONToAttributeValueMap provides a lossy mapping from "basic" json to an AttributeValueMap.
// This allows for the type of "JSON Document" functionality employed in the
// current AWS SDK and outlined in the docs
// (see http://aws.amazon.com/blogs/aws/dynamodb-update-json-and-more/)
func BasicJSONToAttributeValueMap(b []byte) (AttributeValueMap, error) {
	if b == nil {
		return nil, errors.New("arg is nil")
	}
	// unmarshal the arbitrary json
	var i interface{}
	um_err := json.Unmarshal(b, &i)
	if um_err != nil {
		return nil, um_err
	}
	return InterfaceToAttributeValueMap(i)
}

// InterfaceToAttributeValueMap attempts to coerce an appropriate interface {} to
// an AttributeValueMap
func InterfaceToAttributeValueMap(i interface{}) (AttributeValueMap, error) {
	m, m_ok := i.(map[string]interface{})
	if !m_ok {
		return nil, errors.New("top level unmarshal not (map[string] interface{})")
	}
	avm := NewAttributeValueMap()
	for k, v := range m {
		c, cerr := CoerceToAttributeValue(v)
		if cerr != nil {
			return nil, cerr
		}
		avm[k] = c
	}
	return avm, nil
}

// BasicJSONToAttributeValue provides a lossy mapping from "basic" json to an AttributeValue.
// This allows for the type of "JSON Document" functionality employed in the
// current AWS SDK and outlined in the docs
// (see http://aws.amazon.com/blogs/aws/dynamodb-update-json-and-more/)
func BasicJSONToAttributeValue(b []byte) (*AttributeValue, error) {
	if b == nil {
		return nil, errors.New("arg is nil")
	}
	// unmarshal the arbitrary json
	var i interface{}
	um_err := json.Unmarshal(b, &i)
	if um_err != nil {
		return nil, um_err
	}
	return InterfaceToAttributeValue(i)
}

// InterfaceToAttributeValue attempts to coerce an appropriate interface {} to
// an *AttributeValue
func InterfaceToAttributeValue(i interface{}) (*AttributeValue, error) {
	return CoerceToAttributeValue(i)
}

// CoerceToAttributeValue is a lossy translation for basic json to the AWS serialization format
// for AttributeValue. There are types that will be dropped as they are indistinguishable
// without their type designations:
// 1. binary will be dropped as the values will always be coerced to string.
// 2. null (as a type, not a value) will always be coerced to bool.
func CoerceToAttributeValue(i interface{}) (*AttributeValue, error) {
	a := NewAttributeValue()

	// bool (null also coerced to bool)
	b, b_ok := i.(bool)
	if b_ok {
		a.BOOL = new(bool)
		*a.BOOL = b
		return a, nil
	}

	// number - float (the default unmarshal will always use this type)
	n, n_ok := i.(float64)
	if n_ok {
		a.N = strconv.FormatFloat(n, 'f', -1, 64)
		return a, nil
	}

	// string (binary also coerced to string)
	s, s_ok := i.(string)
	if s_ok {
		a.S = s
		return a, nil
	}

	// map of string -> *AttributeValue
	m, m_ok := i.(map[string]interface{})
	if m_ok {
		for k, v := range m {
			a_child, a_child_err := CoerceToAttributeValue(v)
			if a_child_err != nil {
				return nil, a_child_err
			}
			a.M[k] = a_child
		}
		return a, nil
	}

	// the only type of list that is inferred by the generic unmarshal is []interface{}.
	// we need to use further type inference to determine if the list can be made into
	// an NS or SS...or is heterogenous and should be turned into an L
	l, l_ok := i.([]interface{})
	if l_ok {
		l_len := len(l)

		// check first if the list is composed strictly of floats or strings. If so,
		// then we can make a NS or SS list
		float_vals := make([]float64, 0)
		string_vals := make([]string, 0)
		for _, u := range l {
			f, f_ok := u.(float64)
			if f_ok {
				float_vals = append(float_vals, f)
			} else {
				s, s_ok := u.(string)
				if s_ok {
					string_vals = append(string_vals, s)
				}
			}
		}
		floats_len := len(float_vals)
		strings_len := len(string_vals)

		// the list is all floats, turn it into an NS
		if (floats_len == l_len) && (strings_len == 0) {
			for _, f := range float_vals {
				ferr := a.InsertNS_float64(f)
				if ferr != nil {
					return nil, ferr
				}
			}
			return a, nil
		}
		// the list is all strings, turn it into an SS
		if (strings_len == l_len) && (floats_len == 0) {
			for _, v := range string_vals {
				_ = a.InsertSS(v)
			}
			return a, nil
		}

		// the list was not just strictly strings or floats
		for _, v := range l {
			a_child, a_child_err := CoerceToAttributeValue(v)
			if a_child_err != nil {
				return nil, a_child_err
			}
			a.L = append(a.L, a_child)
		}
		return a, nil
	}

	e := fmt.Sprintf("no coercion for %v", i)
	return nil, errors.New(e)
}

// ToBasicJSON provides a mapping from an AttributeValueMap to basic json
// This allows for items from dynamo to be printed in a flat fashion if desired.
func (a AttributeValueMap) ToBasicJSON() ([]byte, error) {
	if a == nil {
		return nil, errors.New("nil AttributeValueMap")
	}
	c, cerr := a.ToInterface()
	if cerr != nil {
		return nil, cerr
	}
	b, merr := json.Marshal(c)
	if merr != nil {
		return nil, merr
	} else {
		return b, nil
	}
}

// AttributeValueMapToInterface converts the map into a map of the key names to interface types
// that do not have type designations, and be marshaled into basic json
func (a AttributeValueMap) ToInterface() (interface{}, error) {
	if a == nil {
		return "", errors.New("nil AttributeValueMap")
	}
	m := make(map[string]interface{})
	for k, v := range a {
		c, cerr := v.ToInterface()
		if cerr != nil {
			return nil, cerr
		} else {
			m[k] = c
		}
	}
	return m, nil
}

// ToBasicJSON provides a mapping from an AttributeValue to basic json
// This allows for items from dynamo to be printed in a flat fashion if desired.
func (a *AttributeValue) ToBasicJSON() ([]byte, error) {
	if a == nil {
		return nil, errors.New("nil AttributeValue")
	}
	c, cerr := a.ToInterface()
	if cerr != nil {
		return nil, cerr
	}
	b, merr := json.Marshal(c)
	if merr != nil {
		return nil, merr
	} else {
		return b, nil
	}
}

// AttributeValueToInterface strips the AttributeValue type designations and returns a structure
// that can be marshaled into basic json.
func (a *AttributeValue) ToInterface() (interface{}, error) {
	if a == nil {
		return "", errors.New("nil AttributeValue")
	}
	if a.BOOL != nil {
		return *a.BOOL, nil
	}
	if a.NULL != nil {
		return *a.NULL, nil
	}
	if a.S != "" {
		return a.S, nil
	}
	if a.B != "" {
		return a.B, nil
	}
	if a.N != "" {
		f, ferr := strconv.ParseFloat(a.N, 64)
		if ferr != nil {
			return nil, ferr
		} else {
			return f, nil
		}
	}
	if len(a.SS) != 0 {
		return a.SS, nil
	}
	if len(a.BS) != 0 {
		return a.BS, nil
	}
	ns_len := len(a.NS)
	if ns_len != 0 {
		ns := make([]float64, ns_len)
		for i, n := range a.NS {
			f, ferr := strconv.ParseFloat(n, 64)
			if ferr != nil {
				return nil, ferr
			} else {
				ns[i] = f
			}
		}
		return ns, nil
	}
	l_len := len(a.L)
	if l_len != 0 {
		ls := make([]interface{}, l_len)
		for i, v := range a.L {
			c, cerr := v.ToInterface()
			if cerr != nil {
				return nil, cerr
			} else {
				ls[i] = c
			}
		}
		return ls, nil
	}
	m_len := len(a.M)
	if m_len != 0 {
		m := make(map[string]interface{})
		for k, v := range a.M {
			c, cerr := v.ToInterface()
			if cerr != nil {
				return nil, cerr
			} else {
				m[k] = c
			}
		}
		return m, nil
	}
	e := fmt.Sprintf("no coercion for %v", a)
	return nil, errors.New(e)
}
