// Support for the DynamoDB GetItem endpoint.
//
// example use:
//
// tests/get_item-livestest.go
//
package get_item

import (
	"encoding/json"
	"errors"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/authreq"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/aws_const"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributestoget"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributevalue"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/capacity"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/expressionattributenames"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/item"
)

const (
	ENDPOINT_NAME      = "GetItem"
	JSON_ENDPOINT_NAME = ENDPOINT_NAME + "JSON"
	GETITEM_ENDPOINT   = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
)

type GetItem struct {
	AttributesToGet          attributestoget.AttributesToGet                   `json:",omitempty"`
	ConsistentRead           bool                                              // false is sane default
	ExpressionAttributeNames expressionattributenames.ExpressionAttributeNames `json:",omitempty"`
	Key                      item.Key
	ProjectionExpression     string `json:",omitempty"`
	ReturnConsumedCapacity   string `json:",omitempty"`
	TableName                string
}

func NewGetItem() *GetItem {
	g := new(GetItem)
	g.Key = item.NewKey()
	g.ExpressionAttributeNames = expressionattributenames.NewExpressionAttributeNames()
	g.AttributesToGet = make(attributestoget.AttributesToGet, 0)
	return g
}

// Get is an alias for backwards compatibility
type Get GetItem

func NewGet() *Get {
	get_item := NewGetItem()
	get := Get(*get_item)
	return &get
}

type Request GetItem

type Response struct {
	Item             item.Item
	ConsumedCapacity *capacity.ConsumedCapacity `json:",omitempty"`
}

type response Response

type response_no_capacity struct {
	Item item.Item
}

func NewResponse() *Response {
	r := new(Response)
	r.Item = item.NewItem()
	r.ConsumedCapacity = capacity.NewConsumedCapacity()
	return r
}

// Some work required to omit empty ConsumedCapacity fields
func (r Response) MarshalJSON() ([]byte, error) {
	if r.ConsumedCapacity.Empty() {
		var ri response_no_capacity
		ri.Item = r.Item
		return json.Marshal(ri)
	}
	ri := response(r)
	return json.Marshal(ri)
}

// ResponseItemJSON can be formed from a Response when the caller wishes
// to receive the Item as basic JSON.
type ResponseItemJSON struct {
	Item             interface{}
	ConsumedCapacity *capacity.ConsumedCapacity `json:",omitempty"`
}

type responseItemJSON ResponseItemJSON

type responseItemJSON_no_capacity struct {
	Item interface{}
}

func NewResponseItemJSON() *ResponseItemJSON {
	r := new(ResponseItemJSON)
	r.ConsumedCapacity = capacity.NewConsumedCapacity()
	return r
}

// Some work required to omit empty ConsumedCapacity fields
func (r ResponseItemJSON) MarshalJSON() ([]byte, error) {
	if r.ConsumedCapacity.Empty() {
		var ri responseItemJSON_no_capacity
		ri.Item = r.Item
		return json.Marshal(ri)
	}
	ri := responseItemJSON(r)
	return json.Marshal(ri)
}

// ToResponseItemJSON will try to convert the Response to a ResponseItemJSON,
// where the interface value for Item represents a structure that can be
// marshaled into basic JSON.
func (resp *Response) ToResponseItemJSON() (*ResponseItemJSON, error) {
	if resp == nil {
		return nil, errors.New("receiver is nil")
	}
	a := attributevalue.AttributeValueMap(resp.Item)
	c, cerr := a.ToInterface()
	if cerr != nil {
		return nil, cerr
	}
	resp_json := NewResponseItemJSON()
	resp_json.ConsumedCapacity = resp.ConsumedCapacity
	resp_json.Item = c
	return resp_json, nil
}

func (get_item *GetItem) EndpointReq() ([]byte, int, error) {
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(get_item)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON, GETITEM_ENDPOINT)
}

func (get *Get) EndpointReq() ([]byte, int, error) {
	get_item := GetItem(*get)
	return get_item.EndpointReq()
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	get_item := GetItem(*req)
	return get_item.EndpointReq()
}
