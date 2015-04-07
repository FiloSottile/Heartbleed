package expected

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Roundtrip some examples
func TestExpectedMarshal(t *testing.T) {
	s := []string{
		`{"MyConstraint1":{"AttributeValueList":[{"S":"a string"}],"ComparisonOperator":"BEGINS_WITH","Value":{"N":"4"},"Exists":true}}`,
		`{"MyConstraint2":{"AttributeValueList":[{"S":"a string"}],"ComparisonOperator":"BEGINS_WITH","Value":{"N":"4"}}}`,
	}
	for _, v := range s {
		var a Expected
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
