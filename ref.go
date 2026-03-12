package golib

import (
	"fmt"
	"math"
)

func Int64RefFromFloat64(v float64) *int64 {
	intV, fractV := math.Modf(v)
	if fractV == 0 {
		return new(int64(intV))
	}
	return nil
}

// Int64OrString checks if data is an int or any of the accepted strings.
func Int64OrString(data interface{}, accepted ...string) (*int64, string, error) {
	if data == nil {
		return nil, "", nil
	}

	switch v := data.(type) {
	case string:
		if StrInArray(v, accepted) {
			return nil, v, nil
		}
	case float64:
		i := Int64RefFromFloat64(v)
		if i != nil {
			return i, "", nil
		}
	}
	return nil, "", fmt.Errorf("Int64OrString: %v Accepted: %v", data, accepted)
}
