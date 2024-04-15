package golib

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplacer1(t *testing.T) {
	r := Replacer{
		EmptyKeys: []*regexp.Regexp{
			regexp.MustCompile(`%.*?%`),
		},
	}
	r.Set(`%a%`, `A`)
	r.Set(`%b%`, `B`)
	if !assert.Equal(t, "Wer A sagt kann auch B sagen aber nicht ", r.Replace("%wegdamit%Wer %a% sagt kann auch %b% sagen aber nicht %c%")) {
		return
	}
}

func TestReplacer2(t *testing.T) {
	r := Replacer{
		repl: map[string]string{
			"%_global_object_id%": "1@dadfbcce-a55d-40ec-baaf-afbe37941304",
			"%_system_object_id%": "1",
			"%object._id%":        "1",
			"%object._version%":   "1",
			"%object.ref%":        "obj1",
		},
		EmptyKeys: []*regexp.Regexp{
			regexp.MustCompile(`%object\.titel%`),
		},
	}
	if !assert.Equal(t, "obj1 ", r.Replace("%object.ref% %object.titel%")) {
		return
	}
}

func TestIntOrReplace(t *testing.T) {
	r := Replacer{}

	r.Set("%min%", "400")

	i64, err := r.IntOrReplace("%min%")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, int64(400), i64) {
		return
	}

	i64, err = r.IntOrReplace(int32(-1234))
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, int64(-1234), i64) {
		return
	}

	// Default path (0,nil)
	i64, err = r.IntOrReplace(nil)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, int64(0), i64) {
		return
	}

}
