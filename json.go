package golib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// Write the given object into file
func JsonWriteFile(fn string, v interface{}) error {
	j, err := json.MarshalIndent(&v, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fn, j, 0644)
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

func JsonUnmarshalReader(r io.Reader, v interface{}) (err error) {
	err = json.NewDecoder(r).Decode(&v)
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
	defer r.Close()
	err := JsonUnmarshalReader(r, v)
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

func JsonBytesIndent(v interface{}, prefix, indent string) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent(prefix, indent)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	bs := buf.Bytes()
	// remove the \n which Encode adds at the end
	return bs[0 : len(bs)-1], nil
}

func JsonBytes(v interface{}) ([]byte, error) {
	return JsonBytesIndent(v, "", "")
}

func JsonString(v interface{}) string {
	return JsonStringIndent(v, "", "    ")
}

func JsonStringIndent(v interface{}, prefix, indent string) string {
	if v == nil {
		return ""
	}
	bs, err := JsonBytesIndent(v, prefix, indent)
	if err != nil {
		return fmt.Sprintf("%#v [Json Error:%s]", v, err)
	}
	return string(bs)
}

// JsonUnmarshalObject marshals the source into json and unmarshals it into target
func JsonUnmarshalObject(source any, target any) error {
	data, err := json.Marshal(source)
	if err != nil {

		return err
	}
	return JsonUnmarshalWithNumber(data, &target)
}

// JsonUnmarshalWithNumber unmarshals data into value. If the json
// decoding returns an error with "cannot unmarshal number ... into float64"
// we try again to unmarshal using "UseNumber"
var numberMatcher = regexp.MustCompile("cannot unmarshal number .*? into Go value of type float64")

func JsonUnmarshalWithNumber(data []byte, value any) (err error) {
	err = json.Unmarshal(data, &value)
	if err != nil {
		if numberMatcher.MatchString(err.Error()) {
			// try again using a json.Number decoder
			dec := json.NewDecoder(bytes.NewReader(data))
			dec.UseNumber()
			return dec.Decode(&value)
		} else {
			return err
		}
	}
	return nil
}

// JsonUnmarshalQuery unmarshals a query string into target Every query
// parameter needs to be unmarshalled separately Otherwise they'll be considered
// json strings and unable to unmarshal in top struct In the case of raw strings
// (no object/array/number etc), it will not unmarshall For this one, we just
// directly assign the raw value
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
