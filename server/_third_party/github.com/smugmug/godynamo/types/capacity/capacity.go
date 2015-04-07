// Package Capacity implements the Capacity type. See:
// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Capacity.html
package capacity

import (
	"encoding/json"
)

type ConsumedCapacityUnit float32

type ConsumedCapacityUnit_struct struct {
	CapacityUnits ConsumedCapacityUnit `json:",omitempty"`
}

type ConsumedCapacity struct {
	CapacityUnits          ConsumedCapacityUnit                   `json:",omitempty"`
	GlobalSecondaryIndexes map[string]ConsumedCapacityUnit_struct `json:",omitempty"`
	LocalSecondaryIndexes  map[string]ConsumedCapacityUnit_struct `json:",omitempty"`
	Table                  *ConsumedCapacityUnit_struct           `json:",omitempty"`
	TableName              string                                 `json:",omitempty"`
}

func NewConsumedCapacity() *ConsumedCapacity {
	c := new(ConsumedCapacity)
	c.GlobalSecondaryIndexes = make(map[string]ConsumedCapacityUnit_struct)
	c.LocalSecondaryIndexes = make(map[string]ConsumedCapacityUnit_struct)
	c.Table = new(ConsumedCapacityUnit_struct)
	return c
}

type ReturnConsumedCapacity string

type consumedcapacity ConsumedCapacity

// Empty determines if this has struct has been assigned
func (c *ConsumedCapacity) Empty() bool {
	if c == nil {
		return true
	}
	if ((c.Table == nil) || (c.Table.CapacityUnits == 0)) &&
		len(c.LocalSecondaryIndexes) == 0 &&
		len(c.GlobalSecondaryIndexes) == 0 &&
		c.TableName == "" &&
		c.CapacityUnits == 0 {
		return true
	} else {
		return false
	}
}

func (c ConsumedCapacity) MarshalJSON() ([]byte, error) {
	if c.Empty() {
		return json.Marshal(nil)
	}
	ci := consumedcapacity(c)
	if c.Table != nil && c.Table.CapacityUnits == 0 {
		ci.Table = nil
	}
	return json.Marshal(ci)
}
