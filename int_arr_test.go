package golib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	ia := IntArr{5, 2, 3, 4}
	ia.Sort()
	if !assert.Equal(t, ia, IntArr{2, 3, 4, 5}) {
		return
	}
}
