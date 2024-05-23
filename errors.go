package golib

// JsonUnmarshalError

type ErrJsonUnmarshal struct {
	err    error
	params struct {
		Value      string `json:"value"`
		SourceType string `json:"sourcetype"`
		TargetType string `json:"targettype"`
		Detail     string `json:"detail"`
	}
}

// JsonUnmarshalError returns a new instance of JsonUnmarshalError with default values
func JsonUnmarshalError() *ErrJsonUnmarshal {
	e := &ErrJsonUnmarshal{}
	return e
}

func (e ErrJsonUnmarshal) ErrorCode() string {
	return "JsonUnmarshalError"
}

func (e ErrJsonUnmarshal) Package() string {
	return "golib"
}

func (e ErrJsonUnmarshal) Error() string {
	return e.params.Detail
}

// Detail sets the value and returns a copy of the error (use for chaining)
func (e *ErrJsonUnmarshal) Detail(v string) *ErrJsonUnmarshal {
	e.params.Detail = v
	return e
}

// GetDetail returns the original error detail
func (e *ErrJsonUnmarshal) GetDetail() string {
	return e.params.Detail
}

// SourceType sets the value and returns a copy of the error (use for chaining)
func (e *ErrJsonUnmarshal) SourceType(v string) *ErrJsonUnmarshal {
	e.params.SourceType = v
	return e
}

// GetSourceType returns the source type
func (e *ErrJsonUnmarshal) GetSourceType() string {
	return e.params.SourceType
}

// TargetType sets the value and returns a copy of the error (use for chaining)
func (e *ErrJsonUnmarshal) TargetType(v string) *ErrJsonUnmarshal {
	e.params.TargetType = v
	return e
}

// GetTargetType returns the target type
func (e *ErrJsonUnmarshal) GetTargetType() string {
	return e.params.TargetType
}

// Value sets the value and returns a copy of the error (use for chaining)
func (e *ErrJsonUnmarshal) Value(v string) *ErrJsonUnmarshal {
	e.params.Value = v
	return e
}

// GetValue returns the original value
func (e *ErrJsonUnmarshal) GetValue() string {
	return e.params.Value
}

// Params returns all parameters as map
func (e ErrJsonUnmarshal) Params() interface{} {
	return e.params
}
