package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt64Intersect(t *testing.T) {
	aOnly, bOnly, both := Int64Intersect([]int64{1, 2, 3}, []int64{2, 3, 4})
	if !assert.ElementsMatch(t, []int64{1}, aOnly) {
		return
	}
	if !assert.ElementsMatch(t, []int64{4}, bOnly) {
		return
	}
	if !assert.ElementsMatch(t, []int64{2, 3}, both) {
		return
	}

	aOnly, bOnly, both = Int64Intersect([]int64{1, 2, 3}, []int64{1, 2, 3})
	if !assert.ElementsMatch(t, []int64{}, aOnly) {
		return
	}
	if !assert.ElementsMatch(t, []int64{}, bOnly) {
		return
	}
	if !assert.ElementsMatch(t, []int64{1, 2, 3}, both) {
		return
	}

	aOnly, bOnly, both = Int64Intersect([]int64{1, 2, 3}, []int64{4, 5, 6})
	if !assert.ElementsMatch(t, []int64{1, 2, 3}, aOnly) {
		return
	}
	if !assert.ElementsMatch(t, []int64{4, 5, 6}, bOnly) {
		return
	}
	if !assert.ElementsMatch(t, []int64{}, both) {
		return
	}

	return
}
