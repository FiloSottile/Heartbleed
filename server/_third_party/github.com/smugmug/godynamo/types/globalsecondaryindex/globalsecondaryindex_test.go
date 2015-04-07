package globalsecondaryindex

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Roundtrip some examples
func TestGlobalSecondaryIndexMarshal(t *testing.T) {
	s := []string{`{"IndexName":"LastPostIndex","KeySchema":[{"AttributeName": "ForumName","KeyType":"HASH"},{"AttributeName":"LastPostDateTime","KeyType": "RANGE"}],"Projection": {"ProjectionType": "KEYS_ONLY"},"ProvisionedThroughput": {"ReadCapacityUnits":200,"WriteCapacityUnits":200}}`}
	for _, v := range s {
		var a GlobalSecondaryIndex
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
