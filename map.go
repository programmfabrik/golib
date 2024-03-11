package golib

// Keys returns all map keys
func MapKeys[M ~map[K]V, K comparable, V any](src M) (keys []K) {
	if src == nil {
		return nil
	}
	keys = make([]K, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	return keys
}
