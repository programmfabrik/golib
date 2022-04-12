package golib

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplacer1(t *testing.T) {
	r := Replacer{
		EmptyKey: regexp.MustCompile(`%.*?%`),
	}
	r.Set(`%a%`, `A`)
	r.Set(`%b%`, `B`)
	if !assert.Equal(t, "Wer A sagt kann auch B sagen aber nicht ", r.Replace("%wegdamit%Wer %a% sagt kann auch %b% sagen aber nicht %c%")) {
		return
	}
}
