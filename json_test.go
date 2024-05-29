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
