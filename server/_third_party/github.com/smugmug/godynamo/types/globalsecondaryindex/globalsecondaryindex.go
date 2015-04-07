// GlobalSecondaryIndex for defining new keys. See:
// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_GlobalSecondaryIndex.html
package globalsecondaryindex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/aws_strings"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/keydefinition"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/provisionedthroughput"
)

type GlobalSecondaryIndex struct {
	IndexName  string
	KeySchema  keydefinition.KeySchema
	Projection struct {
		NonKeyAttributes []string
		ProjectionType   string
	}
	ProvisionedThroughput provisionedthroughput.ProvisionedThroughput
}

func NewGlobalSecondaryIndex() *GlobalSecondaryIndex {
	g := new(GlobalSecondaryIndex)
	g.KeySchema = make(keydefinition.KeySchema, 0)
	g.Projection.NonKeyAttributes = make([]string, 0)
	return g
}

type globalSecondaryIndex GlobalSecondaryIndex

type GlobalSecondaryIndexes []GlobalSecondaryIndex

func (g GlobalSecondaryIndex) MarshalJSON() ([]byte, error) {
	if !(aws_strings.ALL == g.Projection.ProjectionType ||
		aws_strings.KEYS_ONLY == g.Projection.ProjectionType ||
		aws_strings.INCLUDE == g.Projection.ProjectionType) {
		e := fmt.Sprintf("endpoint.GlobalSecondaryIndex.MarshalJSON: "+
			"ProjectionType %s is not valid", g.Projection.ProjectionType)
		return nil, errors.New(e)
	}
	if len(g.Projection.NonKeyAttributes) > 20 {
		e := fmt.Sprintf("endpoint.GlobalSecondaryIndex.MarshalJSON: " +
			"NonKeyAttributes > 20")
		return nil, errors.New(e)
	}
	var gi globalSecondaryIndex
	gi.IndexName = g.IndexName
	gi.KeySchema = g.KeySchema
	gi.Projection = g.Projection
	// if present, must have length between 1 and 20
	if len(g.Projection.NonKeyAttributes) == 0 {
		gi.Projection.NonKeyAttributes = nil
	}
	gi.ProvisionedThroughput = g.ProvisionedThroughput
	gi.Projection.ProjectionType = g.Projection.ProjectionType
	return json.Marshal(gi)
}

type GlobalSecondaryIndexDesc struct {
	IndexName      string
	IndexSizeBytes uint64
	IndexStatus    string
	ItemCount      uint64
	KeySchema      keydefinition.KeySchema
	Projection     struct {
		NonKeyAttributes []string
		ProjectionType   string
	}
	ProvisionedThroughput provisionedthroughput.ProvisionedThroughputDesc
}

func NewGlobalSecondaryIndexDesc() *GlobalSecondaryIndexDesc {
	d := new(GlobalSecondaryIndexDesc)
	d.KeySchema = make(keydefinition.KeySchema, 0)
	d.Projection.NonKeyAttributes = make([]string, 0)
	return d
}

type GlobalSecondaryIndexUpdates struct {
	IndexName             string
	ProvisionedThroughput provisionedthroughput.ProvisionedThroughput
}

func NewGlobalSecondaryIndexUpdates() *GlobalSecondaryIndexUpdates {
	g := new(GlobalSecondaryIndexUpdates)
	return g
}
