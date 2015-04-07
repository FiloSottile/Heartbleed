// Package ItemCollectionMetrics implements the profiling response type. See:
// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_ItemCollectionMetrics.html
package itemcollectionmetrics

import (
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributevalue"
)

type ItemCollectionMetrics struct {
	ItemCollectionKey   attributevalue.AttributeValueMap
	SizeEstimateRangeGB [2]uint64
}

func NewItemCollectionMetrics() *ItemCollectionMetrics {
	i := new(ItemCollectionMetrics)
	i.ItemCollectionKey = attributevalue.NewAttributeValueMap()
	return i
}

// used by BatchWriteItem
type ItemCollectionMetricsMap map[string][]*ItemCollectionMetrics

func NewItemCollectionMetricsMap() ItemCollectionMetricsMap {
	i := make(map[string][]*ItemCollectionMetrics)
	return i
}

type ReturnItemCollectionMetrics string
