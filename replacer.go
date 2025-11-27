package golib

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Global regexp cache for empty key replacers. If used in hot paths the
// compilation of the regexp uses a lot of CPU.
var precReMtx = sync.Mutex{}
var precRe = map[string]*regexp.Regexp{}

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

// AddEmptyKeyReplacer adds all regexps as empty key replacers. It replaces
// matching parts of the string with "" when Replace is called. The regexp are
// cached in this package and never expire.
func (rep *Replacer) AddEmptyKeyReplacer(regexps ...string) {
	precReMtx.Lock()
	defer precReMtx.Unlock()
	for _, s := range regexps {
		if _, has := precRe[s]; !has {
			precRe[s] = regexp.MustCompile(s)
		}
		rep.EmptyKeys = append(rep.EmptyKeys, precRe[s])
	}
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
	switch vt := v.(type) {
	case int:
		return int64(vt), nil
	case uint:
		return int64(vt), nil
	case int8:
		return int64(vt), nil
	case uint8:
		return int64(vt), nil
	case int16:
		return int64(vt), nil
	case uint16:
		return int64(vt), nil
	case int32:
		return int64(vt), nil
	case uint32:
		return int64(vt), nil
	case int64:
		return vt, nil
	case uint64:
		return int64(vt), nil
	case float32:
		return int64(vt), nil
	case float64:
		return int64(vt), nil
	case string:
		// v isn't a string, most likely it's int
		return strconv.ParseInt(rep.Replace(vt), 10, 64)
	default: // also <nil>
		return 0, nil
	}
}

func (rep Replacer) Dump() {
	for _, s := range rep.Debug() {
		Pln(s)
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

	for _, ek := range rep.EmptyKeys {
		ss = append(ss, "empty regexp: "+ek.String())
	}
	return ss
}

func (rep Replacer) Repl() map[string]string {
	return rep.repl
}
