package golib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonUnmarshalQuery(t *testing.T) {
	type customType struct {
		Henk  string `json:"henk"`
		Horst int    `json:"horst"`
	}
	type UploadParamsMultiple struct {
		ID        int64  `json:"id"`
		IDs       string `json:"ids"`
		Slice     []string
		Float64   float64     `json:"float64"`
		Float32   float32     `json:"float32,omitempty"`
		Any       any         `json:"any"`
		StringPtr *string     `json:"str*"`
		Bool      bool        `json:"bool"`
		Custom    *customType `json:"mytype"`
	}
	qv := url.Values{}
	qv.Set("id", "12")
	qv.Set("ids", "12,45")
	qv.Set("float64", "12.34")
	qv.Set("float32", "12.34")
	qv.Set("any", "12.34")
	qv.Set("str*", "strüng")
	qv.Set("bool", "1")

	sl := []string{"1", "2", "3", "4"}
	slBs, _ := json.Marshal(sl)
	qv.Set("Slice", string(slBs))

	mt := customType{
		Henk:  "henk",
		Horst: 123,
	}

	mtBs, _ := json.Marshal(mt)
	qv.Set("mytype", string(mtBs))

	upm := UploadParamsMultiple{}
	err := JsonUnmarshalQuery(qv, &upm)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, int64(12), upm.ID) {
		return
	}
	if !assert.Equal(t, "12,45", qv.Get("ids")) {
		return
	}
	if !assert.Equal(t, sl, upm.Slice) {
		return
	}
	if !assert.Equal(t, float64(12.34), upm.Float64) {
		return
	}
	if !assert.Equal(t, float32(12.34), upm.Float32) {
		return
	}
	if !assert.Equal(t, 12.34, upm.Any) { // parsed to float
		return
	}
	if !assert.Equal(t, true, upm.StringPtr != nil) {
		return
	}
	if !assert.Equal(t, "strüng", *upm.StringPtr) {
		return
	}
	if !assert.Equal(t, true, upm.Bool) {
		return
	}
	if !assert.Equal(t, true, upm.Custom != nil) {
		return
	}
	if !assert.Equal(t, mt, *upm.Custom) {
		return
	}
}

func TestJsonUnmarshalErrorWithPropertyName(t *testing.T) {
	type JsonValues struct {
		Text      string      `json:"text"`
		Array     []string    `json:"array"`
		SubObject *JsonValues `json:"obj"`
	}

	var (
		// target is a struct with properties, the error must include the property names
		target JsonValues
		jue    JsonUnmarshalError
	)

	type testCase struct {
		rawJson  string
		source   string
		target   string
		property string
	}
	for idx, c := range []testCase{
		{
			rawJson:  `{"text": false}`,
			source:   "bool",
			target:   "string",
			property: "JsonValues.text",
		},
		{
			rawJson:  `{"array": false}`,
			source:   "bool",
			target:   "[]string",
			property: "JsonValues.array",
		},
		{
			rawJson:  `{"obj": false}`,
			source:   "bool",
			target:   "golib.JsonValues",
			property: "JsonValues.obj",
		},
		{
			rawJson:  `{"obj": {"array": "invalid"}}`,
			source:   "string",
			target:   "[]string",
			property: "JsonValues.obj.array",
		},
	} {
		err := JsonUnmarshal([]byte(c.rawJson), &target)
		if !assert.Error(t, err) {
			return
		}

		switch {
		case errors.As(err, &jue):
		default:
			t.Errorf("expect JsonUnmarshalError")
			return
		}
		if !assert.Equal(t, c.source, jue.SourceType(), fmt.Sprintf("test case %d: %v", idx, c.rawJson)) {
			return
		}
		if !assert.Equal(t, c.target, jue.TargetType(), fmt.Sprintf("test case %d: %v", idx, c.rawJson)) {
			return
		}
		if !assert.Equal(t, c.property, jue.TargetPropertyName(), fmt.Sprintf("test case %d: %v", idx, c.rawJson)) {
			return
		}
	}
}

func TestJsonUnmarshalErrorWithoutPropertyName(t *testing.T) {
	var (
		// target is a primitive type, not a struct with properties
		target string
		jue    JsonUnmarshalError
	)

	type testCase struct {
		rawJson string
		source  string
		target  string
	}
	for idx, c := range []testCase{
		{
			rawJson: `false`,
			source:  "bool",
			target:  "string",
		},
		{
			rawJson: `[1,2,3]`,
			source:  "array",
			target:  "string",
		},
	} {
		err := JsonUnmarshal([]byte(c.rawJson), &target)
		if !assert.Error(t, err) {
			return
		}

		switch {
		case errors.As(err, &jue):
		default:
			t.Errorf("expect JsonUnmarshalError")
			return
		}
		if !assert.Equal(t, c.source, jue.SourceType(), fmt.Sprintf("test case %d: %v", idx, c.rawJson)) {
			return
		}
		if !assert.Equal(t, c.target, jue.TargetType(), fmt.Sprintf("test case %d: %v", idx, c.rawJson)) {
			return
		}
		// there is no property name, check for empty string
		if !assert.Equal(t, "", jue.TargetPropertyName(), fmt.Sprintf("test case %d: %v", idx, c.rawJson)) {
			return
		}
	}
}
