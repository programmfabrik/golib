package golib

import (
	"encoding/json"
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
