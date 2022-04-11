package golib

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestSetInStruct(t *testing.T) {
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

	ct := cfgTest{}
	err := SetInStruct(map[string]string{
		"INT":              "4",
		"BOOL":             "true",
		"SIMPLE":           "test",
		"INNER_test":       "test",
		"INNER_TESTARR":    `["test1", "test2"]`,
		"INNER_NESTED_DSN": "henk-db",
		// "INNER_MAPME_torsten_NAME": "mein name is torsten",
	}, "_", func(s string) string { return strings.ToUpper(s) }, &ct)
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
	// if !assert.Equal(t, "mein name ist torsten", ct.Inner.MapMe["torsten"].Name) {
	// 	return
	// }
	if !assert.Equal(t, 4, ct.Int) {
		return
	}
}
