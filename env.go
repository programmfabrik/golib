package golib

import (
	"encoding/json"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/yudai/pp"
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
	for key, value := range setMap {
		err = setData(strings.Split(key, sep), value, data, eq)
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

func setData(keyParts []string, value string, data interface{}, eq func(string) string, path ...string) (err error) {
	// pp.Println(keyParts)

	rv := reflect.ValueOf(data).Elem()

	// println("setData", rv.Type().String())
	var rvMap, rvMapIndex reflect.Value
	var haveMap bool
	switch rv.Kind() {
	case reflect.Map:
		rvMap = rv
		rvMapIndex = reflect.ValueOf(keyParts[0])
		rv = rv.MapIndex(rvMapIndex)
		haveMap = true
		if !rvMap.CanAddr() {
			return errNotAddressable
		}
		if !rv.CanAddr() {
			rv2 := rv
			println("unable to addr struct!", rv2.Addr().CanAddr())
		}
		keyParts = keyParts[1:]
	case reflect.Struct:
		if !rv.CanAddr() {
			return errNotAddressable
		}
	default:
		return errNoStruct
	}
	if rv.Kind() != reflect.Struct {
		return errNoStruct
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
			println(">", strings.Join(keyParts[1:], "."), fv.Type().Kind().String(), value)
			if fv.Type().Kind() == reflect.Map {
				if len(keyParts[1:]) < 2 || fv.Type().Key().Kind() != reflect.String {
					return errNoMapString
				}
				// set keyParts[1:2] inside the map
				println(fv.IsNil())
				if fv.IsNil() {
					fv.Set(reflect.MakeMap(fv.Type()))
				}
				mKey := reflect.ValueOf(keyParts[1:2][0])
				fv2 := fv.MapIndex(mKey)
				if !fv2.IsValid() {
					fv.SetMapIndex(mKey, reflect.Zero(fv.Type().Elem()))
				}
				Pln("we have a map! %s v %v", fv.Type().Key(), fv.MapIndex(mKey).String())
			}
			err = setData(keyParts[1:], value, fv.Addr().Interface(), eq, path2...)
			if err != nil {
				return err
			}
		} else {
			kpath := strings.Join(path2, ".")
			println("setting", strings.Join(path2, "."), fv.Type().String(), value)
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
				println("setting value", value)
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
	if haveMap {
		println(rvMap.Type().String())
		pp.Println(rv.Interface())
		rvMap.SetMapIndex(rvMapIndex, rv)
	}

	return nil
}
