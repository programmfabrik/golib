package golib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SliceApply applies all functions Fs against all items in t and returns t2
func TestSliceApply(t *testing.T) {
	s := []string{"a  ", "B", "  c"}
	if !assert.Equal(t, []string{"A", "B", "C"}, SliceApply(s, strings.ToUpper, strings.TrimSpace)) {
		return
	}
	s = []string{"   a  ", "  B  ", "  c  "}
	if !assert.Equal(t, []string{"   a  ", "  b  ", "  c  "}, SliceApply(s, strings.ToLower)) {
		return
	}
	s = []string{"   a  ", "  B  ", "  c  "}
	if !assert.Equal(t, []string{"   a  ", "  B  ", "  c  "}, SliceApply(s)) {
		return
	}
	s = nil
	var s2 []string
	if !assert.Equal(t, s2, SliceApply(s)) {
		return
	}
}
