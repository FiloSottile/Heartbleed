package get_item

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRequestMarshal(t *testing.T) {
	s := []string{
		`{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"HowdoIupdatemultipleitems?"}},"AttributesToGet":["LastPostDateTime","Message","Tags"],"ConsistentRead":true,"ReturnConsumedCapacity":"TOTAL"}`, `{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"HowdoIupdatemultipleitems?"}},"ProjectionExpression":"LastPostDateTime,Message,Tags","ConsistentRead":true,"ReturnConsumedCapacity":"TOTAL"}`,
	}
	for _, v := range s {
		var g GetItem
		um_err := json.Unmarshal([]byte(v), &g)
		if um_err != nil {
			t.Errorf("cannot unmarshal to create:\n" + v + "\n")
		}
		json, jerr := json.Marshal(g)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n", v, string(json))
	}
}

func TestResponseMarshal(t *testing.T) {
	s := []string{`{"ConsumedCapacity":{"CapacityUnits":1,"TableName":"Thread"},"Item":{"Tags":{"SS":["Update","MultipleItems","HelpMe"]},"LastPostDateTime":{"S":"201303190436"},"Message":{"S":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?"}}}`}
	for _, v := range s {
		var g Response
		um_err := json.Unmarshal([]byte(v), &g)
		if um_err != nil {
			t.Errorf("cannot unmarshal to create:\n" + v + "\n")
		}
		json1, jerr := json.Marshal(g)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n", v, string(json1))
		c, cerr := g.ToResponseItemJSON()
		if cerr != nil {
			e := fmt.Sprintf("cannot convert %v\n", cerr)
			t.Errorf(e)
		}
		json2, jerr2 := json.Marshal(c)
		if jerr2 != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("JSON: IN:%v, OUT:%v\n", v, string(json2))
	}
}
