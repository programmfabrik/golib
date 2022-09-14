package golib

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Replacer struct {
	EmptyKeys []*regexp.Regexp // regexp to match empty key
	repl      map[string]string
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

func (rep *Replacer) AddEmptyKeyReplacer(s string) {
	rep.EmptyKeys = append(rep.EmptyKeys, regexp.MustCompile(s))
}

func (rep *Replacer) Replace(s string) string {
	if rep == nil {
		return s
	}
	for key, value := range rep.repl {
		s = strings.ReplaceAll(s, key, value)
	}
	for _, emptyKey := range rep.EmptyKeys {
		s = string(emptyKey.ReplaceAllFunc([]byte(s), func(m []byte) []byte { return nil }))
	}
	return s
}

func (rep Replacer) Dump() {
	for _, s := range rep.Debug() {
		println(s)
	}
}

func (rep Replacer) Debug() (ss []string) {
	keys := []string{}
	for key := range rep.repl {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	ss = make([]string, len(keys))
	for idx, key := range keys {
		ss[idx] = key + "=" + rep.repl[key]
	}
	return ss
}
