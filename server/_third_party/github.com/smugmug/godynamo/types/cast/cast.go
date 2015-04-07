// Functions for casting Go primitive types into AWS types, typically stringified.
package cast

import (
	"encoding/base64"
	"strconv"
)

// AWSParseFLoats normalizes numbers-as-strings for transport.
func AWSParseFloat(s string) (string, error) {
	f, ferr := strconv.ParseFloat(s, 64)
	if ferr != nil {
		return "", ferr
	}
	// aws accepts 38 decimal points, may have to change -1 to that
	return strconv.FormatFloat(f, 'f', -1, 64), nil
}

// AWSParseBinary can test if a string has already been encoded.
func AWSParseBinary(s string) error {
	_, err := base64.StdEncoding.DecodeString(s)
	return err
}
