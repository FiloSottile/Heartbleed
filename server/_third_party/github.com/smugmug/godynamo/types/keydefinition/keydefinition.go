// KeyDefinition is how table keys are described.
package keydefinition

type KeyDefinition struct {
	AttributeName string
	KeyType       string
}

type KeySchema []KeyDefinition
