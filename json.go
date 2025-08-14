package golib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
	"strconv"
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

// countingReader wraps an io.Reader and counts how many bytes have been read.
type countingReader struct {
	r io.Reader
	n int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	readBytes, err := c.r.Read(p)
	c.n += int64(readBytes)
	return readBytes, err
}

// Reads the given file into the given value If the filename ends in ".gz", it
// is opened using the gzip library. Returns the bytes read (unzipped)
func JsonReadFile(fn string, v interface{}) (int64, error) {
	var reader io.Reader
	fh, err := os.Open(fn)
	if err != nil {
		return 0, err
	}
	defer fh.Close()

	if strings.HasSuffix(fn, ".gz") {
		gzip, err := gzip.NewReader(fh)
		if err != nil {
			return 0, err
		}
		defer gzip.Close()
		reader = gzip
	} else {
		reader = fh
	}

	// Wrap the chosen reader in a countingReader.
	cr := &countingReader{r: reader}

	return cr.n, JsonUnmarshalReader(cr, v)
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
	return JsonUnmarshal(data, &target)
}

// JsonUnmarshalQuery unmarshals a query string into target. Only the
// first value of each query key is used. The target is a struct where
// the json tags are used to find the query key. The value is parsed
// according to the type of the target struct field.
func JsonUnmarshalQuery(qv url.Values, target any) (err error) {
	// Ensure that target is a pointer to a struct
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("target must be a non-nil pointer to a struct")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}
	rt := rv.Type()

	sourceData := map[string]any{}
	for i := 0; i < rv.NumField(); i++ {
		// fieldValue := rv.Field(i)
		fInfo := rt.Field(i)

		if !fInfo.IsExported() {
			continue
		}
		parts := strings.Split(fInfo.Tag.Get("json"), ",")
		var fieldName string
		switch len(parts) {
		case 0:
			fieldName = fInfo.Name
		default:
			fieldName = parts[0]
		}
		switch fieldName {
		case "-":
			// skip this
			continue
		case "":
			fieldName = fInfo.Name
		default:
			// already set
		}

		v := qv.Get(fieldName)
		if v == "" {
			continue
		}

		fType := fInfo.Type
		for fType.Kind() == reflect.Pointer {
			fType = fType.Elem()
		}

		// Pln("field %q type %q kind %q name %q value %q", fInfo.Name, fInfo.Type, fType.Kind(), fieldName, v)

		switch fType.Kind() {
		case reflect.String:
			sourceData[fieldName] = v
		case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("JsonUnmarshalQuery: value %q for field %q is not %s", v, fInfo.Name, fInfo.Type.Kind())
			}
			sourceData[fieldName] = i
		case reflect.Float64, reflect.Float32:
			fl, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("JsonUnmarshalQuery: value %q for field %q is not %s", v, fInfo.Name, fInfo.Type.Kind())
			}
			sourceData[fieldName] = fl
		case reflect.Slice, reflect.Interface:
			// value must be JSON
			sourceData[fieldName] = json.RawMessage([]byte(v))
		case reflect.Bool:
			sourceData[fieldName] = GetBool(v)
		default:
			// assume JSON
			sourceData[fieldName] = json.RawMessage([]byte(v))
		}
	}
	err = JsonUnmarshalObject(sourceData, target)
	if err != nil {
		return fmt.Errorf("JsonUnmarshalQuery: %w", err)
	}
	return nil
}
