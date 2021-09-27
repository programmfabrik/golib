package golib

import "regexp"

// RegexpMatch uses regexp to match the first key found in m.
// It returns the key and value found.
func RegexpMatch(regex string, m map[string]string) (retK, retV string) {
	re := regexp.MustCompile(regex)
	for k, v := range m {
		if re.Match([]byte(k)) {
			return k, v
		}
	}
	return "", ""
}
