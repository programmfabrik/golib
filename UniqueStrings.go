package golib

// UniqueStrings returns all unique strings
func UniqueStrings(in []string) (out []string) {
	sMap := map[string]bool{}
	for _, s := range in {
		if !sMap[s] {
			out = append(out, s)
		}
		sMap[s] = true
	}
	return out
}
