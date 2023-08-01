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

// GetEnv returns the os environment. Skips
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
) (valuesSet []string, err error) {
	valuesSet = []string{}
	dataV := reflect.ValueOf(data)
	for key, value := range setMap {
		err = setData(strings.Split(key, sep), value, dataV, eq, &valuesSet)
		if err != nil {
			switch err {
			case errNoStruct:
				return valuesSet, errors.Errorf(`%q needs to be struct but is "%T"`, key, value)
			case errNoMapString:
				return valuesSet, errors.Errorf(`%q needs to be a map[string]`, key)
			case errNotAddressable:
				return valuesSet, errors.Errorf(`%q needs to be addressable`, key)
			default:
				return valuesSet, errors.Wrap(err, "SetInStruct")
			}
		}
	}
	return valuesSet, nil
}

var (
	errNoStruct       = errors.New("No struct")
	errNoMapString    = errors.New("No map string")
	errNotAddressable = errors.New("Not addressable")
)

func setData(keyParts []string, value string, rv reflect.Value, eq func(string) string, valuesSet *[]string, path ...string) (err error) {

	// dereference pointers until we have
	// an element. Initialize nil pointers
	// along the way.
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	if !rv.CanAddr() {
		return fmt.Errorf(`%v needs to be addressable`, path)
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
		} else if !mapElem.CanAddr() {
			// copy mapElem
			elemType := rv.Type().Elem()
			mapElemOld := reflect.Indirect(mapElem)

			mapElem = reflect.New(elemType).Elem()
			if elemType.Kind() == reflect.Pointer {
				mapElem.Set(reflect.New(mapElem.Type().Elem()))
			}
			for i := 0; i < reflect.Indirect(mapElem).NumField(); i++ {
				if reflect.Indirect(mapElem).Field(i).CanSet() {
					reflect.Indirect(mapElem).Field(i).Set(mapElemOld.Field(i))
				}
			}
		}
		path = append(path, keyParts[0])
		keyParts = keyParts[1:]
		origMap = rv
		rv = reflect.Indirect(mapElem)
		// sp = strings.Repeat("  ", len(path))
		// Pln("%s setData %v...%v value: %q kind: %q canAddr: %t", sp, path, keyParts, value, rv.Kind(), rv.CanAddr())
	}
	// Pln("accessing map key %s %s", mapKey, rv.Kind())
	if rv.Kind() == reflect.Struct {
		t := rv.Type()
		matched := false
		if len(keyParts) == 0 {
			return fmt.Errorf(`%v is missing a struct member`, path)
		}
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
				err = setData(keyParts[1:], value, fv, eq, valuesSet, path2...)
				if err != nil {
					return err
				}
			} else {
				// Pln(sp+" %v %s [%s]: %s", path, field.Name, fv.Type().String(), value)
				err = setValue(fv, value)
				if err != nil {
					return errors.Wrapf(err, "Path: %q", strings.Join(path2, "."))
				}
				*valuesSet = append(*valuesSet, strings.Join(append(path, field.Name), "."))
				// thats the leaf of the branch -> set the value
			}
		}
		if !matched {
			// println("SetInStruct: Field not matched", strings.Join(keyParts, "."))
			// currently ignored
		}
	} else {
		err = setValue(rv, value)
		if err != nil {
			return errors.Wrapf(err, "Path: %q", strings.Join(path, "."))
		}
		*valuesSet = append(*valuesSet, strings.Join(path, "."))
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

func setValue(rv reflect.Value, value string) (err error) {
	switch rv.Type().String() {
	case "bool":
		rv.SetBool(GetBool(value))
	case "int", "int64", "int32":
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "Unable to unmarshal %q", value)
		}
		rv.SetInt(i)
	case "string":
		rv.SetString(value)
	case "[]string":
		err = json.Unmarshal([]byte(value), rv.Addr().Interface())
		if err != nil {
			return fmt.Errorf("Unable to unmarshal %q", value)
		}
	default:
		return fmt.Errorf("Unsupported type %q", rv.Type())
	}
	return nil
}
