package attributevalue

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Round trip some data. Note the sets have repeated elements...make sure they
// are eliminated
func TestAttributeValueMarshal(t *testing.T) {
	s := []string{
		`{"S":"a string"}`,
		`{"B":"aGkgdGhlcmUK"}`,
		`{"N":"5"}`,
		`{"BOOL":true}`,
		`{"NULL":false}`,
		`{"SS":["a","b","c","c","c"]}`,
		`{"BS":["aGkgdGhlcmUK","aG93ZHk=","d2VsbCBoZWxsbyB0aGVyZQ=="]}`,
		`{"NS":["42","1","0","0","1","1","1","42"]}`,
		`{"L":[{"S":"a string"},{"L":[{"S":"another string"}]}]}`,
		`{"M":{"key1":{"S":"a string"},"key2":{"L":[{"NS":["42","42","1"]},{"S":"a string"},{"L":[{"S":"another string"}]}]}}}`,
	}
	for _, v := range s {
		_ = fmt.Sprintf("--------\n")
		_ = fmt.Sprintf("IN:%v\n", v)
		var a AttributeValue
		um_err := json.Unmarshal([]byte(v), &a)
		if um_err != nil {
			_ = fmt.Sprintf("%v\n", um_err)
			t.Errorf("cannot unmarshal\n")
		}

		json, jerr := json.Marshal(a)
		if jerr != nil {
			_ = fmt.Sprintf("%v\n", jerr)
			t.Errorf("cannot marshal\n")
			return
		}
		_ = fmt.Sprintf("OUT:%v\n", string(json))
	}
}

// Demonstrate the use of the Valid function
func TestAttributeValueInvalid(t *testing.T) {
	a := NewAttributeValue()
	a.N = "1"
	a.S = "a"
	if a.Valid() {
		_, jerr := json.Marshal(a)
		if jerr == nil {
			t.Errorf("should not have been able to marshal\n")
			return
		} else {
			_ = fmt.Sprintf("%v\n", jerr)
		}
	}
	a = NewAttributeValue()
	a.N = "1"
	a.B = "fsdfa"
	if a.Valid() {
		_, jerr := json.Marshal(a)
		if jerr == nil {
			t.Errorf("should not have been able to marshal\n")
			return
		} else {
			_ = fmt.Sprintf("%v\n", jerr)
		}
	}

	a = NewAttributeValue()
	a.N = "1"
	a.InsertSS("a")
	if a.Valid() {
		_, jerr := json.Marshal(a)
		if jerr == nil {
			t.Errorf("should not have been able to marshal\n")
			return
		} else {
			_ = fmt.Sprintf("%v\n", jerr)
		}
	}
}

// Empty AttributeValues should emit null
func TestAttributeValueEmpty(t *testing.T) {
	a := NewAttributeValue()
	json_bytes, jerr := json.Marshal(a)
	if jerr != nil {
		_ = fmt.Sprintf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	_ = fmt.Sprintf("OUT:%v\n", string(json_bytes))

	var a2 AttributeValue
	json_bytes, jerr = json.Marshal(a2)
	if jerr != nil {
		_ = fmt.Sprintf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	_ = fmt.Sprintf("OUT:%v\n", string(json_bytes))
}

// Test the Insert funtions
func TestAttributeValueInserts(t *testing.T) {
	a1 := NewAttributeValue()
	a1.InsertSS("hi")
	a1.InsertSS("hi") // duplicate, should be removed
	a1.InsertSS("bye")
	json, jerr := json.Marshal(a1)
	if jerr != nil {
		_ = fmt.Sprintf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	_ = fmt.Sprintf("OUT:%v\n", string(json))
}

// Test the Insert functions
func TestAttributeValueInserts2(t *testing.T) {
	a1 := NewAttributeValue()
	_ = a1.InsertSS("hi")
	_ = a1.InsertSS("hi") // duplicate, should be removed
	_ = a1.InsertSS("bye")
	json_bytes, jerr := json.Marshal(a1)
	if jerr != nil {
		_ = fmt.Sprintf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	_ = fmt.Sprintf("OUT:%v\n", string(json_bytes))

	a2 := NewAttributeValue()
	_ = a2.InsertL(a1)
	a1 = nil // should be fine, above line should make a new copy
	json_bytes, jerr = json.Marshal(a2)
	if jerr != nil {
		_ = fmt.Sprintf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	_ = fmt.Sprintf("OUT:%v\n", string(json_bytes))

	a3 := NewAttributeValue()
	nerr := a3.InsertN("fred")
	if nerr == nil {
		t.Errorf("should have returned error from InsertN\n")
		return
	} else {
		_ = fmt.Sprintf("%v\n", nerr)
	}
	berr := a3.InsertB("1")
	if berr == nil {
		t.Errorf("should have returned error from InsertB\n")
		return
	} else {
		_ = fmt.Sprintf("%v\n", berr)
	}
}

// Should fail, a2 is uninitialized
func TestBadCopy(t *testing.T) {
	a1 := NewAttributeValue()
	_ = a1.InsertSS("hi")
	_ = a1.InsertSS("bye")

	var a2 = new(AttributeValue)

	cp_err := a1.Copy(a2)
	if a2 == nil {
		t.Errorf("should have returned error from Copy\n")
		return
	} else {
		_ = fmt.Sprintf("%v\n", cp_err)
	}
}

// Make sure Valid emits as null
func TestAttributeValueUpdate(t *testing.T) {
	a := NewAttributeValueUpdate()
	a.Action = "DELETE"
	json_bytes, jerr := json.Marshal(a)
	if jerr != nil {
		_ = fmt.Sprintf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	_ = fmt.Sprintf("OUT:%v\n", string(json_bytes))

}

func TestCoerceAttributeValueBasicJSON(t *testing.T) {
	js := []string{`{"a":{"b":"c"},"d":[{"e":"f"},"g","h"],"i":[1.0,2.0,3.0],"j":["x","y"]}`,
		`"a"`, `true`,
		`[1,2,3,2,3]`}
	for _, i := range js {
		_ = fmt.Sprintf("--------\n")
		j := []byte(i)
		_ = fmt.Sprintf("IN:%v\n", string(j))
		av, av_err := BasicJSONToAttributeValue(j)
		if av_err != nil {
			_ = fmt.Sprintf("%v\n", av_err)
			t.Errorf("cannot coerce")
			return
		}
		av_json, av_json_err := json.Marshal(av)
		if av_json_err != nil {
			_ = fmt.Sprintf("%v\n", av_json_err)
			t.Errorf("cannot marshal")
			return
		}
		_ = fmt.Sprintf("OUT:%v\n", string(av_json))
		b, cerr := av.ToBasicJSON()
		if cerr != nil {
			_ = fmt.Sprintf("%v\n", cerr)
			t.Errorf("cannot coerce")
			return
		}
		_ = fmt.Sprintf("RT:%v\n", string(b))
	}
}

func TestCoerceAttributeValueMapBasicJSON(t *testing.T) {
	js := []string{`{"AS":"1234string","AN":3,"ANS":[1,2,1,2,3],"ASS":["a","a","b"],"ABOOL":true,"AL":["1234string",3,[1,2,3],["a","b"]],"AM":{"AMS":"1234string","AMN":3,"AMNS":[1,2,3],"AMSS":["a","b"],"AMBOOL":true,"AL":["1234string",3,[1,2,3],["a","b"]]}}`}
	for _, i := range js {
		_ = fmt.Sprintf("--------\n")
		j := []byte(i)
		_ = fmt.Sprintf("IN:%v\n", string(j))
		av, av_err := BasicJSONToAttributeValueMap(j)
		if av_err != nil {
			_ = fmt.Sprintf("%v\n", av_err)
			t.Errorf("cannot coerce")
			return
		}
		av_json, av_json_err := json.Marshal(av)
		if av_json_err != nil {
			_ = fmt.Sprintf("%v\n", av_json_err)
			t.Errorf("cannot marshal")
			return
		}
		_ = fmt.Sprintf("OUT:%v\n", string(av_json))
		b, cerr := av.ToBasicJSON()
		if cerr != nil {
			_ = fmt.Sprintf("%v\n", cerr)
			t.Errorf("cannot coerce")
			return
		}
		_ = fmt.Sprintf("RT:%v\n", string(b))

	}
}
