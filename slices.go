package golib

func SliceFilter[T any](t []T, keep func(T) bool) (t2 []T) {
	if keep == nil {
		return t
	}
	for _, item := range t {
		if keep(item) {
			t2 = append(t2, item)
		}
	}
	return t2
}
