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

// JsonUnmarshal marshals the source into json and unmarshals it into target.
// If there is an error, it checks for known json parser errors and
// if there is a match, a JsonUnmarshalError with parsed information is returned.
// The error messages are formatted by https://pkg.go.dev/encoding/json#Unmarshal
// * "json: cannot unmarshal <value> into Go value of type <type>"
// * "json: cannot unmarshal <value> into Go struct field <target property name> of type <type>"
func JsonUnmarshal(source []byte, target any) (err error) {
	err = json.Unmarshal(source, target)
	if err == nil {
		return nil
	}
	regex := regexp.MustCompile(`json: cannot unmarshal ([^\s]+) into Go value of type (.+)$`)
	matches := regex.FindStringSubmatch(err.Error())
	if len(matches) == 3 {
		return NewJsonUnmarshalError(
			err,
			matches[1], // source type
			matches[2], // target type
			"",         // no target property name
		)
	}

	regex = regexp.MustCompile(`json: cannot unmarshal ([^\s]+) into Go struct field ([^\s]+?) of type (.+)$`)
	matches = regex.FindStringSubmatch(err.Error())
	if len(matches) == 4 {
		return NewJsonUnmarshalError(
			err,
			matches[1], // source type
			matches[3], // target type
			matches[2], // target property name
		)
	}

	// no regex match, just return the original error
	return err
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
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return dec.Decode(&value)
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
