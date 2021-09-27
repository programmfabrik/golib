package golib

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

// Write the given object into file
func JsonWriteFile(fn string, v interface{}) error {
	j, err := json.MarshalIndent(&v, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fn, j, 0644)
	return err
}

// Reads the given file into the given value
// If the filename ends in ".gz", it is opened
// using the gzip library.
func JsonReadFile(fn string, v interface{}) error {
	var reader io.Reader
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer fh.Close()

	if strings.HasSuffix(fn, ".gz") {
		gzip, err := gzip.NewReader(fh)
		if err != nil {
			return err
		}
		defer gzip.Close()
		reader = gzip
	} else {
		reader = fh
	}

	return JsonUnmarshalReader(reader, v)
}

func JsonUnmarshalReader(r io.Reader, v interface{}) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&v)
	if err != nil {
		return fmt.Errorf("json.Decode io.reader error: %w", err)
	}
	return nil
}

func JsonUnmarshalReaderStrict(r io.Reader, v interface{}) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	err := dec.Decode(&v)
	if err != nil {
		return fmt.Errorf("json.Decode io.reader error: %w", err)
	}
	return nil
}

func JsonUnmarshalReadCloser(r io.ReadCloser, v interface{}) error {
	err := JsonUnmarshalReader(r, v)
	defer r.Close()
	if err != nil {
		return err
	}
	return nil
}

func JsonPretty(b []byte) string {
	var data interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return fmt.Sprintf("%q [Json Error:%s]", string(b), err)
	}
	b2, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Sprintf("%q [Json Error:%s]", string(b), err)
	}
	return string(b2)
}

func JsonString(v interface{}) string {
	if v == nil {
		return ""
	}
	bytes, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		// fall back to the built in
		return fmt.Sprintf("%#v [Json Error:%s]", v, err)
	}
	return string(bytes)
}

// JsonUnmarshalObject marshals the source into json and unmarshals it into target
func JsonUnmarshalObject(source interface{}, target interface{}) error {
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &target)
}

// JsonUnmarshalQuery unmarshals a query string into target
// Every query parameter needs to be unmarshalled separately
// Otherwise they'll be considered json strings and unable to unmarshal in top struct
// In the case of raw strings (no object/array/number etc), it will not unmarshall
// For this one, we just directly assign the raw value
func JsonUnmarshalQuery(qv url.Values, target interface{}) error {
	m := map[string]interface{}{}
	for k, vs := range qv {
		var mm interface{}
		err := json.Unmarshal([]byte(vs[0]), &mm)
		if err != nil {
			mm = vs[0]
		}
		switch v := mm.(type) {
		case string:
			if v == "" {
				continue
			}
			m[k] = v
		default:
			m[k] = v
		}
	}
	return JsonUnmarshalObject(m, target)
}
