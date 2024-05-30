package golib

type JsonUnmarshalError struct {
	err                error // original error for reference
	sourceType         string
	targetType         string
	targetPropertyName string
}

// NewJsonUnmarshalError returns a new JsonUnmarshalError
// The original error is stored together with the source type and target type
// If no property name is available it is given as an empty string
func NewJsonUnmarshalError(err error, sourceType, targetType, targetPropertyName string) JsonUnmarshalError {
	return JsonUnmarshalError{
		err:                err,
		sourceType:         sourceType,
		targetType:         targetType,
		targetPropertyName: targetPropertyName,
	}
}

// Error returns the original error message
func (jue JsonUnmarshalError) Error() string {
	return jue.err.Error()
}

// SourceType returns the type which could not be unmarshaled
func (jue *JsonUnmarshalError) SourceType() string {
	return jue.sourceType
}

// TargetType returns the type the source type could not be unmarshaled into
func (jue *JsonUnmarshalError) TargetType() string {
	return jue.targetType
}

// TargetPropertyName returns the name of the property the source type could not be unmarshaled into
func (jue *JsonUnmarshalError) TargetPropertyName() string {
	return jue.targetPropertyName
}
