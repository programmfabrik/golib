package lib

import "sort"

type StringMatcher map[string]bool

func (sm *StringMatcher) Append(ss ...string) {
	for _, s := range ss {
		(*sm)[s] = true
	}
}

func (sm StringMatcher) Match(s string) bool {
	_, ok := sm[s]
	return ok
}

func (sm StringMatcher) SortedKeys() (sorted []string) {
	sorted = make([]string, 0, len(sm))
	for k := range sm {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	return sorted
}
