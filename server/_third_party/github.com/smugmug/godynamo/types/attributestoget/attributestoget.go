// AttributesToGet is used in GetItem etc
package attributestoget

type AttributesToGet []string

type attributesToGet AttributesToGet

// AttributesToGet is already a reference type
func NewAttributesToGet() AttributesToGet {
	s := make([]string, 0)
	return s
}
