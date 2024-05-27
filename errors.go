package golib

type ErrJsonUnmarshal struct {
	sourceType string
	targetType string
	detail     string
}

// JsonUnmarshalError returns a new instance of JsonUnmarshalError with default values
func JsonUnmarshalError() *ErrJsonUnmarshal {
	e := &ErrJsonUnmarshal{}
	return e
}

func (e ErrJsonUnmarshal) ErrorCode() string {
	return "JsonUnmarshalError"
}

func (e ErrJsonUnmarshal) Error() string {
	return e.detail
}

// Detail sets the value and returns a copy of the error (use for chaining)
func (e *ErrJsonUnmarshal) Detail(v string) *ErrJsonUnmarshal {
	e.detail = v
	return e
}

// GetDetail returns the original error detail
func (e *ErrJsonUnmarshal) GetDetail() string {
	return e.detail
}

// SourceType sets the value and returns a copy of the error (use for chaining)
func (e *ErrJsonUnmarshal) SourceType(v string) *ErrJsonUnmarshal {
	e.sourceType = v
	return e
}

// GetSourceType returns the source type
func (e *ErrJsonUnmarshal) GetSourceType() string {
	return e.sourceType
}

// TargetType sets the value and returns a copy of the error (use for chaining)
func (e *ErrJsonUnmarshal) TargetType(v string) *ErrJsonUnmarshal {
	e.targetType = v
	return e
}

// GetTargetType returns the target type
func (e *ErrJsonUnmarshal) GetTargetType() string {
	return e.targetType
}
