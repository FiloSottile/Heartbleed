// Support for the DynamoDB PutItem endpoint.
//
// example use:
//
// tests/put_item-livestest.go
//
package put_item

import (
	"encoding/json"
	"errors"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/authreq"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/aws_const"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributesresponse"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributevalue"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/aws_strings"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/expected"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/expressionattributenames"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/item"
)

const (
	ENDPOINT_NAME      = "PutItem"
	JSON_ENDPOINT_NAME = ENDPOINT_NAME + "JSON"
	PUTITEM_ENDPOINT   = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	// the permitted ReturnValues flags for this op
	RETVAL_ALL_OLD = aws_strings.RETVAL_ALL_OLD
	RETVAL_NONE    = aws_strings.RETVAL_NONE
)

type PutItem struct {
	ConditionExpression         string                                            `json:",omitempty"`
	ConditionalOperator         string                                            `json:",omitempty"`
	Expected                    expected.Expected                                 `json:",omitempty"`
	ExpressionAttributeNames    expressionattributenames.ExpressionAttributeNames `json:",omitempty"`
	ExpressionAttributeValues   attributevalue.AttributeValueMap                  `json:",omitempty"`
	Item                        item.Item
	ReturnConsumedCapacity      string `json:",omitempty"`
	ReturnItemCollectionMetrics string `json:",omitempty"`
	ReturnValues                string `json:",omitempty"`
	TableName                   string
}

// NewPut will return a pointer to an initialized PutItem struct.
func NewPutItem() *PutItem {
	p := new(PutItem)
	p.Expected = expected.NewExpected()
	p.ExpressionAttributeNames = expressionattributenames.NewExpressionAttributeNames()
	p.ExpressionAttributeValues = attributevalue.NewAttributeValueMap()
	p.Item = item.NewItem()
	return p
}

type Request PutItem

// Put is an alias for backwards compatibility
type Put PutItem

func NewPut() *Put {
	put_item := NewPutItem()
	put := Put(*put_item)
	return &put
}

// PutItemJSON differs from PutItem in that JSON is a string, which allows you to use a basic
// JSON document as the Item
type PutItemJSON struct {
	ConditionExpression         string                                            `json:",omitempty"`
	ConditionalOperator         string                                            `json:",omitempty"`
	Expected                    expected.Expected                                 `json:",omitempty"`
	ExpressionAttributeNames    expressionattributenames.ExpressionAttributeNames `json:",omitempty"`
	ExpressionAttributeValues   attributevalue.AttributeValueMap                  `json:",omitempty"`
	Item                        interface{}
	ReturnConsumedCapacity      string `json:",omitempty"`
	ReturnItemCollectionMetrics string `json:",omitempty"`
	ReturnValues                string `json:",omitempty"`
	TableName                   string
}

// NewPutJSON will return a pointer to an initialized PutItemJSON struct.
func NewPutItemJSON() *PutItemJSON {
	p := new(PutItemJSON)
	p.Expected = expected.NewExpected()
	p.ExpressionAttributeNames = expressionattributenames.NewExpressionAttributeNames()
	p.ExpressionAttributeValues = attributevalue.NewAttributeValueMap()
	return p
}

// ToPutItem will attempt to convert a PutItemJSON to PutItem
func (put_item_json *PutItemJSON) ToPutItem() (*PutItem, error) {
	if put_item_json == nil {
		return nil, errors.New("receiver is nil")
	}
	a, cerr := attributevalue.InterfaceToAttributeValueMap(put_item_json.Item)
	if cerr != nil {
		return nil, cerr
	}
	p := NewPutItem()
	p.ConditionExpression = put_item_json.ConditionExpression
	p.ConditionalOperator = put_item_json.ConditionalOperator
	p.Expected = put_item_json.Expected
	p.ExpressionAttributeNames = put_item_json.ExpressionAttributeNames
	p.ExpressionAttributeValues = put_item_json.ExpressionAttributeValues
	p.Item = item.Item(a)
	p.ReturnConsumedCapacity = put_item_json.ReturnConsumedCapacity
	p.ReturnItemCollectionMetrics = put_item_json.ReturnItemCollectionMetrics
	p.ReturnValues = put_item_json.ReturnValues
	p.TableName = put_item_json.TableName
	return p, nil
}

type Response attributesresponse.AttributesResponse

func NewResponse() *Response {
	a := attributesresponse.NewAttributesResponse()
	r := Response(*a)
	return &r
}

func (put_item *PutItem) EndpointReq() ([]byte, int, error) {
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(put_item)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON, PUTITEM_ENDPOINT)
}

func (put *Put) EndpointReq() ([]byte, int, error) {
	put_item := PutItem(*put)
	return put_item.EndpointReq()
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	put_item := PutItem(*req)
	return put_item.EndpointReq()
}

// ValidItem validates the size of a json serialization of an Item.
// AWS says items can only be 400k bytes binary
func ValidItem(i string) bool {
	return !(len([]byte(i)) > 400000)
}
