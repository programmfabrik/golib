package golib

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yudai/pp"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("FYLR_DEBUG", "Horst")
	os.Setenv("FYLR_debug_HENK", "Horst2")
	envMap := GetEnv("FYLR_(.*)")
	if !assert.Equal(t, "Horst", envMap["DEBUG"]) {
		return
	}
	if !assert.Equal(t, "Horst2", envMap["debug_HENK"]) {
		return
	}
	envMap = GetEnv("FYLR_.*")
	if !assert.Equal(t, "Horst", envMap["FYLR_DEBUG"]) {
		return
	}
}

func TestMapValues1(t *testing.T) {
	vMap := MapValues([]string{"horst=schröder", "henk=schrader=d"}, "")
	if !assert.Equal(t, "schröder", vMap["horst"]) {
		return
	}
	if !assert.Equal(t, "schrader=d", vMap["henk"]) {
		return
	}
}

func TestMapValues2(t *testing.T) {
	vMap := MapValues([]string{"F_horst=schröder", "F_henk=schrader=d"}, "F_(.*)")
	if !assert.Equal(t, "schröder", vMap["horst"]) {
		return
	}
	if !assert.Equal(t, "schrader=d", vMap["henk"]) {
		return
	}
}

type cfgTest struct {
	Int    int
	Bool   bool
	Simple string
	Inner  struct {
		Test    string
		TestArr []string
		Nested  struct {
			DSN string
		}
		MapMe map[string]struct {
			Name  string
			Value int
		}
	}
}

func WEGTestSetMapValue(t *testing.T) {
	a := "henk"

	rvA := reflect.Indirect(reflect.ValueOf(&a))
	println("rvA", rvA.CanAddr(), rvA.String())
	rvA.SetString("horst")
	println(a)

	c := cfgTest{}

	rv1 := reflect.ValueOf(&c).Elem()
	println("rv1", rv1.CanAddr(), rv1.String())
	rv2 := rv1.FieldByName("Inner")
	println("rv2", rv2.CanAddr(), rv2.Type().String())
	rv3 := rv2.FieldByName("MapMe")
	println("rv3", rv3.CanAddr(), rv3.String(), rv2.Type().String())
	rv3.Set(reflect.MakeMap(rv3.Type()))

	var mapElem reflect.Value
	elemType := rv3.Type().Elem()
	if !mapElem.IsValid() {
		mapElem = reflect.New(elemType).Elem()
	} else {
		mapElem.Set(reflect.Zero(elemType))
	}
	rv3.SetMapIndex(reflect.ValueOf("henk"), mapElem)

	pp.Println(c)
}

func TestSetInStruct(t *testing.T) {

	ct := cfgTest{}
	err := SetInStruct(map[string]string{
		"INT":                      "4",
		"BOOL":                     "true",
		"SIMPLE":                   "test",
		"INNER_test":               "test",
		"INNER_TESTARR":            `["test1", "test2"]`,
		"INNER_NESTED_DSN":         "henk-db",
		"INNER_MAPME_torsten_NAME": "mein name is torsten",
	}, "_",
		func(s string) string {
			return strings.ToUpper(s)
		},
		&ct)

	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, "test", ct.Simple) {
		return
	}
	if !assert.ElementsMatch(t, []string{"test1", "test2"}, ct.Inner.TestArr) {
		return
	}
	if !assert.Equal(t, "henk-db", ct.Inner.Nested.DSN) {
		return
	}
	if !assert.Equal(t, "mein name is torsten", ct.Inner.MapMe["torsten"].Name) {
		return
	}
	if !assert.Equal(t, 4, ct.Int) {
		return
	}
}
