package golib

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonUnmarshalWithNumber(t *testing.T) {
	a := []byte(`{"horst":2.62264311E+82434647}`)
	var n any
	// fails with regular json unmarshal
	err := json.Unmarshal(a, &n)
	if !assert.Error(t, err) {
		return
	}
	// does not fail
	err = JsonUnmarshalWithNumber(a, &n)
	if !assert.NoError(t, err) {
		return
	}
	bs, _ := json.Marshal(n)
	if !assert.Equal(t, string(a), string(bs)) {
		return
	}
}
