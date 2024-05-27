package golib

type ErrJsonUnmarshal struct {
	err                error
	sourceType         string
	targetType         string
	targetPropertyName string
}

func JsonUnmarshalError() *ErrJsonUnmarshal {
	e := &ErrJsonUnmarshal{}
	return e
}

func (e ErrJsonUnmarshal) ErrorCode() string {
	return "JsonUnmarshalError"
}

func (e ErrJsonUnmarshal) Error() string {
	return e.err.Error()
}

func (e *ErrJsonUnmarshal) OriginalError(err error) *ErrJsonUnmarshal {
	e.err = err
	return e
}

func (e *ErrJsonUnmarshal) GetOriginalError() error {
	return e.err
}

func (e *ErrJsonUnmarshal) SourceType(v string) *ErrJsonUnmarshal {
	e.sourceType = v
	return e
}

func (e *ErrJsonUnmarshal) GetSourceType() string {
	return e.sourceType
}

func (e *ErrJsonUnmarshal) TargetType(v string) *ErrJsonUnmarshal {
	e.targetType = v
	return e
}

func (e *ErrJsonUnmarshal) GetTargetType() string {
	return e.targetType
}

func (e *ErrJsonUnmarshal) TargetPropertyName(v string) *ErrJsonUnmarshal {
	e.targetPropertyName = v
	return e
}

func (e *ErrJsonUnmarshal) GetTargetPropertyName() string {
	return e.targetPropertyName
}
