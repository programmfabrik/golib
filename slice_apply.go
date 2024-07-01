package golib

// SliceApply applies all functions Fs against all items in t and returns t2
func SliceApply[T ~[]E, E any](t T, Fs ...func(E) E) (t2 T) {
	if len(Fs) == 0 || len(t) == 0 {
		return t
	}
	for _, item := range t {
		for _, f := range Fs {
			item = f(item)
		}
		t2 = append(t2, item)
	}
	return t2
}
