package itemcollectionmetrics

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestExpectedMarshal(t *testing.T) {
	s := []string{
		`{"ItemCollectionKey":{"AttributeValue":{"S":"a string"}},"SizeEstimateRangeGB":[0,10]}`,
	}
	for _, v := range s {
		var a ItemCollectionMetrics
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
