// A collection of types and methods common to all of the endpoint/* packages.
// Define types for the core datatype listed here:
// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Types.html
// Some of the data types at that url are not defined here, rather they may be
// defined with special limits for the various endpoints - see the packages
// for these endpoints for details.
package endpoint

import (
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributedefinition"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributesresponse"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributestoget"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributevalue"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/aws_strings"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/capacity"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/expected"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/globalsecondaryindex"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/item"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/itemcollectionmetrics"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/keydefinition"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/localsecondaryindex"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/nullable"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/provisionedthroughput"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/returnvalues"
	"net/http"
)

// re-exported consts for backwards compatibility
const (
	HASH               = aws_strings.HASH
	RANGE              = aws_strings.RANGE
	HASH_KEY_ELEMENT   = aws_strings.HASH_KEY_ELEMENT
	RANGE_KEY_ELEMENT  = aws_strings.RANGE_KEY_ELEMENT
	S                  = aws_strings.S
	N                  = aws_strings.N
	B                  = aws_strings.B
	BOOL               = aws_strings.BOOL
	NULL               = aws_strings.NULL
	L                  = aws_strings.L
	M                  = aws_strings.M
	SS                 = aws_strings.SS
	NS                 = aws_strings.NS
	BS                 = aws_strings.BS
	RETVAL_NONE        = aws_strings.RETVAL_NONE
	RETVAL_ALL_OLD     = aws_strings.RETVAL_ALL_OLD
	RETVAL_ALL_NEW     = aws_strings.RETVAL_ALL_NEW
	RETVAL_UPDATED_OLD = aws_strings.RETVAL_UPDATED_OLD
	RETVAL_UPDATED_NEW = aws_strings.RETVAL_UPDATED_NEW
	ALL                = aws_strings.ALL
	SIZE               = aws_strings.SIZE
	TOTAL              = aws_strings.TOTAL
	KEYS_ONLY          = aws_strings.KEYS_ONLY
	INCLUDE            = aws_strings.INCLUDE
	SELECT_ALL         = aws_strings.SELECT_ALL
	SELECT_PROJECTED   = aws_strings.SELECT_PROJECTED
	SELECT_SPECIFIC    = aws_strings.SELECT_SPECIFIC
	SELECT_COUNT       = aws_strings.SELECT_COUNT
	OP_EQ              = aws_strings.OP_EQ
	OP_NE              = aws_strings.OP_NE
	OP_LE              = aws_strings.OP_LE
	OP_LT              = aws_strings.OP_LT
	OP_GE              = aws_strings.OP_GE
	OP_GT              = aws_strings.OP_GT
	OP_NULL            = aws_strings.OP_NULL
	OP_NOT_NULL        = aws_strings.OP_NOT_NULL
	OP_CONTAINS        = aws_strings.OP_CONTAINS
	OP_NOT_CONTAINS    = aws_strings.OP_NOT_CONTAINS
	OP_BEGINS_WITH     = aws_strings.OP_BEGINS_WITH
	OP_IN              = aws_strings.OP_IN
	OP_BETWEEN         = aws_strings.OP_BETWEEN
)

// re-exported types for backwards compatibility
type AttributeDefinition attributedefinition.AttributeDefinition

type AttributeDefinitions attributedefinition.AttributeDefinitions

type AttributeValue attributevalue.AttributeValue

type AttributesResponse attributesresponse.AttributesResponse

type AttributesToGet attributestoget.AttributesToGet

type ConsumedCapacityUnit capacity.ConsumedCapacityUnit

type ConsumedCapacity capacity.ConsumedCapacity

type Item item.Item

type ReturnValues returnvalues.ReturnValues

type NullableString nullable.NullableString

type NullableUInt64 nullable.NullableUInt64

type ReturnConsumedCapacity capacity.ReturnConsumedCapacity

type ReturnItemCollectionMetrics itemcollectionmetrics.ItemCollectionMetrics

type ItemCollectionMetrics itemcollectionmetrics.ItemCollectionMetrics

type Constraints expected.Constraints

type Expected expected.Expected

type ProvionedThroughPut provisionedthroughput.ProvisionedThroughput

type ProvionedThroughPutDesc provisionedthroughput.ProvisionedThroughputDesc

type KeyDefinition keydefinition.KeyDefinition

type KeySchema keydefinition.KeySchema

type LocalSecondaryIndex localsecondaryindex.LocalSecondaryIndex

type LocalSecondaryIndexes localsecondaryindex.LocalSecondaryIndexes

type LocalSecondaryIndexDesc localsecondaryindex.LocalSecondaryIndexDesc

type GlobalSecondaryIndex globalsecondaryindex.GlobalSecondaryIndex

type GlobalSecondaryIndexes globalsecondaryindex.GlobalSecondaryIndexes

type GlobalSecondaryIndexDesc globalsecondaryindex.GlobalSecondaryIndexDesc

// ---------------------------------------------------------------------
// GoDynamo-specific types

// Endpoint is a core interface for defining all endpoints.
// Packages implementing the Endpoint interface should return the
// string output from the authorized request (or ""), the http code,
// and an error (or nil). This is the fundamental endpoint interface of
// GoDynamo.
type Endpoint interface {
	EndpointReq() ([]byte, int, error)
}

// Endpoint_Response describes the response from Dynamo for a given request.
type Endpoint_Response struct {
	Body []byte
	Code int
	Err  error
}

// ReqErr is a convenience function to see if the request was bad.
func ReqErr(code int) bool {
	return code >= http.StatusBadRequest &&
		code < http.StatusInternalServerError
}

// ServerErr is a convenience function to see if the remote server had an
// internal error.
func ServerErr(code int) bool {
	return code >= http.StatusInternalServerError
}

// HttpErr is a convenience function to see determine if the code is an error code.
func HttpErr(code int) bool {
	return ReqErr(code) || ServerErr(code)
}
