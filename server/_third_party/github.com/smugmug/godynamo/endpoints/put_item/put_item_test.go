// Tests JSON formats as described on the AWS docs site. For live tests, see ../../tests
package put_item

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRequestMarshal(t *testing.T) {
	s := []string{`{"TableName":"Thread","Item":{"LastPostDateTime":{"S":"201303190422"},"Tags":{"SS":["Update","MultipleItems","HelpMe"]},"ForumName":{"S":"AmazonDynamoDB"},"Message":{"S":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?"},"Subject":{"S":"HowdoIupdatemultipleitems?"},"LastPostedBy":{"S":"fred@example.com"}},"Expected":{"ForumName":{"Exists":false},"Subject":{"Exists":false}}}`, `{"TableName":"Thread","Item":{"LastPostDateTime":{"S":"201303190422"},"Tags":{"SS":["Update","MultipleItems","HelpMe"]},"ForumName":{"S":"AmazonDynamoDB"},"Message":{"S":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?"},"Subject":{"S":"HowdoIupdatemultipleitems?"},"LastPostedBy":{"S":"fred@example.com"}},"ConditionExpression":"ForumName<>:fandSubject<>:s","ExpressionAttributeValues":{":f":{"S":"AmazonDynamoDB"},":s":{"S":"HowdoIupdatemultipleitems?"}}}`}
	for _, v := range s {
		var p PutItem
		um_err := json.Unmarshal([]byte(v), &p)
		if um_err != nil {
			t.Errorf("cannot unmarshal RequestItems:\n" + v + "\n")
		}
		json, jerr := json.Marshal(p)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n", v, string(json))
	}
}

func TestRequestJSONMarshal(t *testing.T) {
	s := []string{`{"TableName":"Thread","Item":{"LastPostDateTime":"201303190422","Tags":["Update","MultipleItems","HelpMe"],"ForumName":"AmazonDynamoDB","Message":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?","Subject":"HowdoIupdatemultipleitems?","LastPostedBy":"fred@example.com"},"Expected":{"ForumName":{"Exists":false},"Subject":{"Exists":false}}}`, `{"TableName":"Thread","Item":{"LastPostDateTime":"201303190422","Tags":["Update","MultipleItems","HelpMe"],"ForumName":"AmazonDynamoDB","Message":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?","Subject":"HowdoIupdatemultipleitems?","LastPostedBy":"fred@example.com"},"ConditionExpression":"ForumName<>:fandSubject<>:s","ExpressionAttributeValues":{":f":{"S":"AmazonDynamoDB"},":s":{"S":"HowdoIupdatemultipleitems?"}}}`}
	for _, v := range s {
		var p PutItemJSON
		um_err := json.Unmarshal([]byte(v), &p)
		if um_err != nil {
			t.Errorf("cannot unmarshal RequestItems:\n" + v + "\n")
		}
		json, jerr := json.Marshal(p)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("JSON IN:%v, OUT:%v\n", v, string(json))
	}
}

func TestResponseMarshal(t *testing.T) {
	s := []string{`{"Attributes":{"LastPostedBy":{"S":"alice@example.com"},"ForumName":{"S":"AmazonDynamoDB"},"LastPostDateTime":{"S":"20130320010350"},"Tags":{"SS":["Update","MultipleItems","HelpMe"]},"Subject":{"S":"Maximumnumberofitems?"},"Views":{"N":"5"},"Message":{"S":"Iwanttoput10milliondataitemstoanAmazonDynamoDBtable.Isthereanupperlimit?"}}}`}
	for _, v := range s {
		var p Response
		um_err := json.Unmarshal([]byte(v), &p)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n" + v + "\n")
		}
		json, jerr := json.Marshal(p)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n", v, string(json))
	}
}
