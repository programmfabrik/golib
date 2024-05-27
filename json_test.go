package golib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
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

func TestJsonUnmarshalQuery(t *testing.T) {
	type UploadParamsMultiple struct {
		References      []string `json:"references"`
		ProduceVersions []bool   `json:"produce_versions"`
		VersionNames    []string `json:"version_names"`
		IDParents       []int64  `json:"id_parents"`
	}
	q := url.URL{RawQuery: "access_token=HENK&check_for_duplicates=1"}
	u := UploadParamsMultiple{}
	err := JsonUnmarshalQuery(q.Query(), &u)
	if !assert.NoError(t, err) {
		return
	}
}

type JsonValues struct {
	Text      string      `json:"text"`
	Integer   int         `json:"integer"`
	Decimal   float64     `json:"decimal"`
	Bool      bool        `json:"bool"`
	Array     []any       `json:"array"`
	SubObject *JsonValues `json:"obj"`
}

func TestJsonUnmarshalError(t *testing.T) {
	var (
		target          JsonValues
		unmarshalErr    ErrJsonUnmarshal
		unmarshalErrPtr *ErrJsonUnmarshal
	)

	type testCase struct {
		rawJson                    string
		expectedSourceType         string
		expectedTargetType         string
		expectedTargetPropertyName string
	}
	for idx, c := range []testCase{
		{
			rawJson:                    `{"text": false}`,
			expectedSourceType:         "bool",
			expectedTargetType:         "string",
			expectedTargetPropertyName: "JsonValues.text",
		},
		{
			rawJson:                    `{"text": 123}`,
			expectedSourceType:         "number",
			expectedTargetType:         "string",
			expectedTargetPropertyName: "JsonValues.text",
		},
		{
			rawJson:                    `{"text": 123.456}`,
			expectedSourceType:         "number",
			expectedTargetType:         "string",
			expectedTargetPropertyName: "JsonValues.text",
		},
		{
			rawJson:                    `{"integer": "invalid"}`,
			expectedSourceType:         "string",
			expectedTargetType:         "int",
			expectedTargetPropertyName: "JsonValues.integer",
		},
		{
			rawJson:                    `{"decimal": "invalid"}`,
			expectedSourceType:         "string",
			expectedTargetType:         "float64",
			expectedTargetPropertyName: "JsonValues.decimal",
		},
		{
			rawJson:                    `{"bool": 123}`,
			expectedSourceType:         "number",
			expectedTargetType:         "bool",
			expectedTargetPropertyName: "JsonValues.bool",
		},
		{
			rawJson:                    `{"array": false}`,
			expectedSourceType:         "bool",
			expectedTargetType:         "[]interface {}",
			expectedTargetPropertyName: "JsonValues.array",
		},
	} {
		err := JsonUnmarshal([]byte(c.rawJson), &target)
		if !assert.Error(t, err) {
			return
		}

		switch {
		case errors.As(err, &unmarshalErr):
		case errors.As(err, &unmarshalErrPtr):
			unmarshalErr = *unmarshalErrPtr
		default:
			t.Errorf("expect JsonUnmarshalError")
			return
		}
		if !assert.Equal(t,
			c.expectedSourceType,
			unmarshalErr.GetSourceType(),
			fmt.Sprintf("test case %d: %v: check SourceType", idx, c.rawJson),
		) {
			return
		}
		if !assert.Equal(t,
			c.expectedTargetPropertyName,
			unmarshalErr.GetTargetPropertyName(),
			fmt.Sprintf("test case %d: %v: check TargetPropertyName", idx, c.rawJson),
		) {
			return
		}
		if !assert.Equal(t,
			c.expectedTargetType,
			unmarshalErr.GetTargetType(),
			fmt.Sprintf("test case %d: %v: check TargetType", idx, c.rawJson),
		) {
			return
		}
	}
}
