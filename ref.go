package golib

import (
	"fmt"
	"math"
	"time"
)

// StringRef returns a refernce to the given string
func StringRef(s string, def ...*string) *string {
	switch len(def) {
	case 0:
		return &s
	case 1:
		if s == "" {
			return def[0]
		}
		return &s
	default:
		panic("multiple default params")
	}
}

// Int64Ref returns a refernce to the given string
func Int64Ref(i int64, def ...*int64) *int64 {
	switch len(def) {
	case 0:
		return &i
	case 1:
		if i == 0 {
			return def[0]
		}
		return &i
	default:
		panic("multiple default params")
	}
}

// Float64Ref returns a refernce to the given float64
func Float64Ref(f float64, def ...*float64) *float64 {
	switch len(def) {
	case 0:
		return &f
	case 1:
		if f == 0 {
			return def[0]
		}
		return &f
	default:
		panic("multiple default params")
	}
}

// BoolRef
func BoolRef(b bool, def ...*bool) *bool {
	switch len(def) {
	case 0:
		return &b
	case 1:
		if !b {
			return def[0]
		}
		return &b
	default:
		panic("multiple default params")
	}
}

// IntRef returns a refernece to the given string
func IntRef(i int, def ...*int) *int {
	switch len(def) {
	case 0:
		return &i
	case 1:
		if i == 0 {
			return def[0]
		}
		return &i
	default:
		panic("multiple default params")
	}
}

// TimeRef returns a refernece to the given Time
func TimeRef(t time.Time, def ...*time.Time) *time.Time {
	switch len(def) {
	case 0:
		return &t
	case 1:
		if t.IsZero() {
			return def[0]
		}
		return &t
	default:
		panic("multiple default params")
	}
}

// DurationRef returns a refernece to the given Duration
func DurationRef(d time.Duration, def ...*time.Duration) *time.Duration {
	switch len(def) {
	case 0:
		return &d
	case 1:
		if d == 0 {
			return def[0]
		}
	default:
		panic("multiple default params")
	}
	return &d
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
	return nil, "", fmt.Errorf("Int64OrString: %v Accepted: %v", data, accepted)
}
