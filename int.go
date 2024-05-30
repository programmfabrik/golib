package golib

import (
	"errors"
	"math"
)

func Float64ToInt64(f float64) (i int64, err error) {
	iP, fP := math.Modf(f)
	if fP != 0 {
		return 0, errors.New("Float has fractional part")
	}
	return int64(iP), nil
}

// Int64SliceToInterfaceSlice converts a slice of int64 to slice of interface
func Int64SliceToInterfaceSlice(i []int64) []interface{} {
	in := make([]interface{}, len(i))
	for idx, i0 := range i {
		in[idx] = i0
	}
	return in
}

// Int64Merge returns a new slice with all items from a and b
// merged together, duplicates removed
func Int64Merge(a, b []int64) (merged []int64) {
	merged, bOnly, both := Int64Intersect(a, b)
	return append(merged, append(bOnly, both...)...)
}

// Int64Intersect intersects a and b. It returns three slices,
// 1. items only in a, 2. items only in b, 3. items in both.
func Int64Intersect(a, b []int64) (aOnly, bOnly, both []int64) {
	mA, mB, all :=
		map[int64]bool{},
		map[int64]bool{},
		map[int64]bool{}

	for _, aV := range a {
		mA[aV] = true
		all[aV] = true
	}
	for _, bV := range b {
		mB[bV] = true
		all[bV] = true
	}
	both, aOnly, bOnly = []int64{}, []int64{}, []int64{}
	for v := range all {
		switch {
		case mA[v] && mB[v]:
			both = append(both, v)
		case mA[v]:
			aOnly = append(aOnly, v)
		case mB[v]:
			bOnly = append(bOnly, v)
		}
	}
	return aOnly, bOnly, both
}
