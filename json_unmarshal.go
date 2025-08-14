package golib

import (
	"encoding/json"
	"errors"
)

// JsonUnmarshal marshals the source into json and unmarshals it into target.
// If there is an error, it checks for known json parser errors and
// if there is a match, a JsonUnmarshalError with parsed information is returned.
// The error messages are formatted by https://pkg.go.dev/encoding/json#Unmarshal
// * "json: cannot unmarshal <value> into Go value of type <type>"
// * "json: cannot unmarshal <value> into Go struct field <target property name> of type <type>"

func JsonUnmarshal(source []byte, target any) (err error) {
	err = json.Unmarshal(source, target)
	if err == nil {
		return nil
	}
	// Prefer structured error data (works in json/v2 and v1, though fields are what we need here).
	var ute *json.UnmarshalTypeError
	if errors.As(err, &ute) {
		// ute.Value: the JSON value kind/string (e.g., "string", "number", "object"…)
		// ute.Type:  the Go reflect.Type it tried to put it into
		// ute.Field: dotted path to the struct field (empty when the target is a non-struct type)
		return NewJsonUnmarshalError(
			err,
			ute.Value,         // source type (JSON value)
			ute.Type.String(), // target type (Go type)
			ute.Field,         // target property name (may be "")
		)
	}

	// Not a type-mismatch style error; return as-is.
	return err
}
