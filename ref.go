package golib

import (
	"math"
	"time"

	"github.com/pkg/errors"
)

// StringRef returns a refernce to the given string
func StringRef(s string) *string {
	return &s
}

// Int64Ref returns a refernce to the given string
func Int64Ref(i int64) *int64 {
	return &i
}

// Float64Ref returns a refernce to the given float64
func Float64Ref(i float64) *float64 {
	return &i
}

// BoolRef
func BoolRef(b bool) *bool {
	return &b
}

// IntRef returns a refernece to the given string
func IntRef(i int) *int {
	return &i
}

// TimeRef returns a refernece to the given Time
func TimeRef(t time.Time) *time.Time {
	return &t
}

// DurationRef returns a refernece to the given Duration
func DurationRef(t time.Duration) *time.Duration {
	return &t
}

func Int64RefFromFloat64(v float64) *int64 {
	intV, fractV := math.Modf(v)
	if fractV == 0 {
		return Int64Ref(int64(intV))
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
	return nil, "", errors.Errorf("Int64OrString: %v Accepted: %v", data, accepted)
}
