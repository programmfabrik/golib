package golib

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type EnvMap map[string]string

// GetEnv returns the os environement. Skips
// values not matching the prefix regexp re.
func GetEnv(re string) EnvMap {
	return MapValues(os.Environ(), re)
}

// Env reassembles EnvMap to Env
func (em EnvMap) Env() (env []string) {
	// initialize with <nil>
	for k, v := range em {
		env = append(env, k+"="+v)
	}
	return env
}

// MapValues returns the full set of variables optionally starting with matching
// regexp "re"
func MapValues(values []string, re string) (envMap map[string]string) {
	envMap = map[string]string{}
	reMatch := regexp.MustCompile("^(.*?)=(.*)$")
	var reKeyMatch *regexp.Regexp
	if re != "" {
		reKeyMatch = regexp.MustCompile(re)
	}
	for _, env := range values {
		envMatch := reMatch.FindAllStringSubmatch(env, -1)
		if envMatch == nil {
			continue
		}
		if re == "" {
			envMap[envMatch[0][1]] = envMatch[0][2]
			continue
		}
		key := envMatch[0][1]
		keyM := reKeyMatch.FindStringSubmatch(key)
		if keyM == nil {
			continue
		}
		if len(keyM) > 1 {
			envMap[keyM[1]] = envMatch[0][2]
		} else {
			envMap[keyM[0]] = envMatch[0][2]
		}
	}
	return envMap
}

// SetInStruct uses setMap to set value in passed data.
// The keys take the form of PREFIX_KEY
// KEYs are split by sep and compared with struct members using the eq func
// the passed data
func SetInStruct(
	setMap map[string]string, // source of data
	sep string, // key path separator (e.g. "." or "_")
	eq func(string) string, // compare function (last part of path with the struct field name)
	data interface{}, // target to set the data in
) (err error) {
	dataV := reflect.ValueOf(data)
	for key, value := range setMap {
		err = setData(strings.Split(key, sep), value, dataV, eq)
		if err != nil {
			switch err {
			case errNoStruct:
				return errors.Errorf(`%q needs to be struct but is "%T"`, key, value)
			case errNoMapString:
				return errors.Errorf(`%q needs to be a map[string]`, key)
			case errNotAddressable:
				return errors.Errorf(`%q needs to be addressable`, key)
			default:
				return errors.Wrap(err, "SetInStruct")
			}
		}
	}
	return nil
}

var (
	errNoStruct       = errors.New("No struct")
	errNoMapString    = errors.New("No map string")
	errNotAddressable = errors.New("Not addressable")
)

func setData(keyParts []string, value string, rv reflect.Value, eq func(string) string, path ...string) (err error) {

	// dereference pointers until we have
	// an element. Initialize nil pointers
	// along the way.
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	// sp := strings.Repeat("  ", len(path))
	// Pln("%s setData %v...%v value: %q kind: %q canAddr: %t", sp, path, keyParts, value, rv.Kind(), rv.CanAddr())

	var mapKey, mapElem, origMap reflect.Value

	switch rv.Kind() {
	case reflect.Map:
		// Create map if needed
		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rv.Type()))
		}

		// Create map element if needed
		mapKey = reflect.ValueOf(keyParts[0])
		mapElem = rv.MapIndex(mapKey)
		if !mapElem.IsValid() {
			elemType := rv.Type().Elem()
			mapElem = reflect.New(elemType).Elem()
			if elemType.Kind() == reflect.Pointer {
				mapElem.Set(reflect.New(mapElem.Type().Elem()))
			}
			// Pln("eme type %s %v %t", elemType, elemType.Kind() == reflect.Pointer, mapElem.IsValid())
		}
		path = append(path, keyParts[0])
		keyParts = keyParts[1:]
		origMap = rv
		rv = reflect.Indirect(mapElem)
		// sp = strings.Repeat("  ", len(path))
		// Pln("%s setData %v...%v value: %q kind: %q canAddr: %t", sp, path, keyParts, value, rv.Kind(), rv.CanAddr())
	case reflect.Struct:
		if !rv.CanAddr() {
			return fmt.Errorf(`%v needs to be addressable`, path)
		}
	}
	// Pln("accessing map key %s %s", mapKey, rv.Kind())
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf(`%v needs to be struct but is "%T"`, path, value)
	}
	t := rv.Type()
	matched := false
	for i := 0; i < rv.NumField(); i++ {
		field := t.Field(i)
		if eq(field.Name) != eq(keyParts[0]) || !field.IsExported() {
			continue
		}
		matched = true
		fv := rv.Field(i)
		path2 := make([]string, len(path))
		copy(path2, path)
		path2 = append(path2, field.Name)

		if len(keyParts) > 1 {
			// more parts left, dive
			err = setData(keyParts[1:], value, fv, eq, path2...)
			if err != nil {
				return err
			}
		} else {
			kpath := strings.Join(path2, ".")
			// Pln(sp+" %v %s [%s]: %s", path, field.Name, fv.Type().String(), value)
			switch fv.Type().String() {
			case "bool":
				fv.SetBool(GetBool(value))
			case "int", "int64", "int32":
				i, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return errors.Wrapf(err, "Unable to unmarshal %q into key %q", value, kpath)
				}
				fv.SetInt(i)
			case "string":
				fv.SetString(value)
			case "[]string":
				err = json.Unmarshal([]byte(value), fv.Addr().Interface())
				if err != nil {
					return errors.Wrapf(err, "Unable to unmarshal %q into key %q", value, kpath)
				}
			default:
				return errors.Errorf("Unsupported type %q for key %q", t, kpath)
			}
			// thats the leaf of the branch -> set the value
		}
	}
	if !matched {
		println("field not matched")
		// currently ignored
	}
	// If we access an element of a map, set the value, unless
	// it is a pointer.
	if mapElem.IsValid() {
		if mapElem.Kind() == reflect.Pointer {
			origMap.SetMapIndex(mapKey, rv.Addr())
		} else {
			origMap.SetMapIndex(mapKey, rv)
		}
	}
	return nil
}
