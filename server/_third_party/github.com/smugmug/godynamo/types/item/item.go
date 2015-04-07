// Item as described in docs for various endpoints.
package item

import (
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributevalue"
)

type Item attributevalue.AttributeValueMap

type item Item

// Item is already a reference type
func NewItem() Item {
	a := attributevalue.NewAttributeValueMap()
	return Item(a)
}

// ItemLike is an interface for those structs you wish to map back and forth to Items.
// This is currently provided instead of the lossy translation advocated by the
// JSON document mapping as described by AWS.
type ItemLike interface {
	ToItem(interface{}) (Item, error)
	FromItem(Item) (interface{}, error)
}

// GetItem and UpdateItem share a Key type which is another alias to AttributeValueMap
type Key attributevalue.AttributeValueMap

func NewKey() Key {
	a := attributevalue.NewAttributeValueMap()
	return Key(a)
}
