// Support for AttributeDefinition type. See
// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_AttributeDefinition.html
package attributedefinition

// Do not omit these if they are empty, they are both required
type AttributeDefinition struct {
	AttributeName string
	AttributeType string
}

type AttributeDefinitions []AttributeDefinition
