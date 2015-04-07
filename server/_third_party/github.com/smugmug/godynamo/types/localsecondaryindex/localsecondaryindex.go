// LocalSecondaryIndex for defining new keys. See:
// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_LocalSecondaryIndex.html
package localsecondaryindex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/aws_strings"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/keydefinition"
)

type LocalSecondaryIndex struct {
	IndexName  string                  `json:",omitempty"`
	KeySchema  keydefinition.KeySchema `json:",omitempty"`
	Projection struct {
		NonKeyAttributes []string `json:",omitempty"`
		ProjectionType   string   `json:",omitempty"`
	}
}

func NewLocalSecondaryIndex() *LocalSecondaryIndex {
	l := new(LocalSecondaryIndex)
	l.KeySchema = make(keydefinition.KeySchema, 0)
	l.Projection.NonKeyAttributes = make([]string, 0)
	return l
}

type localSecondaryIndex LocalSecondaryIndex

type LocalSecondaryIndexes []LocalSecondaryIndex

func (l LocalSecondaryIndex) MarshalJSON() ([]byte, error) {
	if !(aws_strings.ALL == l.Projection.ProjectionType ||
		aws_strings.KEYS_ONLY == l.Projection.ProjectionType ||
		aws_strings.INCLUDE == l.Projection.ProjectionType) {
		e := fmt.Sprintf("endpoint.LocalSecondaryIndex.MarshalJSON: "+
			"ProjectionType %s is not valid", l.Projection.ProjectionType)
		return nil, errors.New(e)
	}
	if len(l.Projection.NonKeyAttributes) > 20 {
		e := fmt.Sprintf("endpoint.LocalSecondaryIndex.MarshalJSON: " +
			"NonKeyAttributes > 20")
		return nil, errors.New(e)
	}
	var li localSecondaryIndex
	li.IndexName = l.IndexName
	li.KeySchema = l.KeySchema
	li.Projection = l.Projection
	// if present, must have length between 1 and 20
	if l.Projection.NonKeyAttributes != nil && len(l.Projection.NonKeyAttributes) != 0 {
		li.Projection.NonKeyAttributes = l.Projection.NonKeyAttributes
	} else {
		li.Projection.NonKeyAttributes = nil
	}
	li.Projection.ProjectionType = l.Projection.ProjectionType
	return json.Marshal(li)
}

type LocalSecondaryIndexDesc struct {
	IndexName      string
	IndexSizeBytes uint64
	ItemCount      uint64
	KeySchema      keydefinition.KeySchema
	Projection     struct {
		NonKeyAttributes []string
		ProjectionType   string
	}
}

func NewLocalSecondaryIndexDesc() *LocalSecondaryIndexDesc {
	d := new(LocalSecondaryIndexDesc)
	d.KeySchema = make(keydefinition.KeySchema, 0)
	d.Projection.NonKeyAttributes = make([]string, 0)
	return d
}
