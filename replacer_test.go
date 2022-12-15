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
