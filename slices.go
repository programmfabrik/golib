package golib

func SliceFilter[T ~[]E, E any](t T, keep func(E) bool) (t2 T) {
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
