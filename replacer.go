package golib

import (
	"strconv"
	"strings"
)

type Replacer struct {
	repl map[string]string
}

func (rep *Replacer) Set(key, value string) {
	if rep.repl == nil {
		rep.repl = map[string]string{}
	}
	rep.repl[key] = value
}

func (rep *Replacer) SetInt(key string, value int) {
	rep.Set(key, strconv.Itoa(value))
}

func (rep *Replacer) SetInt64(key string, value int64) {
	rep.Set(key, strconv.FormatInt(value, 10))
}

func (rep *Replacer) Replace(s string) string {
	if rep == nil {
		return s
	}
	for key, value := range rep.repl {
		s = strings.ReplaceAll(s, key, value)
	}
	return s
}
