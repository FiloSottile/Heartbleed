// Package Expected implements the ExpectedAttributeValue data type. See:
// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_ExpectedAttributeValue.html
// The documentation defines it in places as Expected and/or ExpectedAttributeValue.
package expected

import (
	"encoding/json"
	"errors"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/types/attributevalue"
)

type Expected map[string]*Constraints

func NewExpected() Expected {
	e := make(Expected)
	return e
}

type Constraints struct {
	AttributeValueList []*attributevalue.AttributeValue
	ComparisonOperator string
	Value              *attributevalue.AttributeValue
	Exists             *bool
}

func NewConstraints() *Constraints {
	c := new(Constraints)
	c.Exists = new(bool)
	c.Value = attributevalue.NewAttributeValue()
	c.AttributeValueList = make([]*attributevalue.AttributeValue, 0)
	return c
}

type constraints Constraints

func (c *Constraints) UnmarshalJSON(data []byte) error {
	if c == nil {
		return errors.New("pointer receiver for unmarshal is nil")
	}
	var ci constraints
	t_err := json.Unmarshal(data, &ci)
	if t_err != nil {
		return t_err
	}

	if c.Exists == nil {
		c.Exists = new(bool)
	}

	if ci.Exists == nil {
		*c.Exists = true
	} else {
		*c.Exists = *ci.Exists
	}

	if ci.Value != nil {
		c.Value = attributevalue.NewAttributeValue()
		cp_err := ci.Value.Copy(c.Value)
		if cp_err != nil {
			return cp_err
		}
	}

	l_ci_avl := len(ci.AttributeValueList)
	if l_ci_avl != 0 {
		c.AttributeValueList = make([]*attributevalue.AttributeValue, l_ci_avl)
		for i, _ := range ci.AttributeValueList {
			c.AttributeValueList[i] = attributevalue.NewAttributeValue()
			cp_err := ci.AttributeValueList[i].Copy(c.AttributeValueList[i])
			if cp_err != nil {
				return cp_err
			}
		}
	}
	return nil
}

func (c Constraints) MarshalJSON() ([]byte, error) {
	var ci constraints
	if c.Exists != nil {
		ci.Exists = new(bool)
		*ci.Exists = *c.Exists
	}
	ci.Value = c.Value
	ci.AttributeValueList = c.AttributeValueList
	ci.ComparisonOperator = c.ComparisonOperator
	return json.Marshal(ci)
}
