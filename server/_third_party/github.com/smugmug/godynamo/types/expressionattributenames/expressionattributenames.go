// ExpressionAttributeNames is used in GetItem etc
package expressionattributenames

type ExpressionAttributeNames map[string]string

func NewExpressionAttributeNames() ExpressionAttributeNames {
	e := make(map[string]string)
	return e
}
