package golib

// UniqueInt64s returns all unique strings
func UniqueInt64s(in []int64) (out []int64) {
	iMap := map[int64]bool{}
	for _, i := range in {
		if iMap[i] {
			out = append(out, i)
			iMap[i] = true
		}
	}
	return out
}
