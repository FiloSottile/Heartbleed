// AttributeResponse is response for PutItem,UpdateItem,DeleteItem,etc
package attributesresponse

import (
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributevalue"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/capacity"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/itemcollectionmetrics"
)

type AttributesResponse struct {
	Attributes            attributevalue.AttributeValueMap             `json:",omitempty"`
	ConsumedCapacity      *capacity.ConsumedCapacity                   `json:",omitempty"`
	ItemCollectionMetrics *itemcollectionmetrics.ItemCollectionMetrics `json:",omitempty"`
}

func NewAttributesResponse() *AttributesResponse {
	a := new(AttributesResponse)
	a.Attributes = attributevalue.NewAttributeValueMap()
	a.ConsumedCapacity = capacity.NewConsumedCapacity()
	a.ItemCollectionMetrics = itemcollectionmetrics.NewItemCollectionMetrics()
	return a
}
