package golib

import "encoding/json"

type NullInt64 struct {
	Value int64
	Valid bool
}

func (i *NullInt64) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.Valid = true

	// leave value at zero if we received a "null"
	if string(data) == "null" {
		return nil
	}

	// The key isn't set to null
	var temp int64
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	i.Value = temp
	return nil
}
