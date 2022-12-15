package golib

import (
	"fmt"
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

// IntOrReplace replaces v if it is a string and parses it as int64. If v
// already is an int, an int64 is returned. If the string (after replacements)
// cannot be parsed as int64, an error is returned.
func (rep *Replacer) IntOrReplace(v any) (int64, error) {
	rs, ok := v.(string)
	if !ok {
		// v isn't a string, most likely it's int
		rs = fmt.Sprintf("%v", v)
	}
	return strconv.ParseInt(rep.Replace(rs), 10, 64)
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
