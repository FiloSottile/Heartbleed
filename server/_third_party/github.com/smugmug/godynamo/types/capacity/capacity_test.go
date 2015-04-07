package capacity

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Roundtrip some examples
func TestCapacityMarshal(t *testing.T) {
	s := []string{
		`{"CapacityUnits":1.00,"TableName":"mytable"}`,
		`{"CapacityUnits":1.01,"TableName":"mytable","Table":{"CapacityUnits":2.01}}`,
		`{"CapacityUnits":1.01,"TableName":"mytable","Table":{"CapacityUnits":2.01},"LocalSecondaryIndexes":{"mylsi":{"CapacityUnits":10.10}},"GlobalSecondaryIndexes":{"mygsi0":{"CapacityUnits":11.11},"mygsi1":{"CapacityUnits":10.10}}}`,
	}
	for _, v := range s {
		var a ConsumedCapacity
		um_err := json.Unmarshal([]byte(v), &a)
		if um_err != nil {
			_ = fmt.Sprintf("%v\n", um_err)
			t.Errorf("cannot unmarshal\n")
		}

		json, jerr := json.Marshal(a)
		if jerr != nil {
			_ = fmt.Sprintf("%v\n", jerr)
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n", v, string(json))
	}
}
