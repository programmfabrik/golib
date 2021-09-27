package golib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueFilename(t *testing.T) {
	efs := []string{
		"/henk/was.jpg",
		"/henk/was.jpg",
		"/henk/i am unique",
		"/henk/youKnowMe",
		"/henk/youKNOWMe",
	}

	uniqMap := UniqueFilename{
		"/henk/was.jpg": true,
	}

	efs2 := []string{}
	for _, ef := range efs {
		efs2 = append(efs2, uniqMap.Add(ef))
	}

	if !assert.Equal(t, []string{
		"/henk/was_00001.jpg",
		"/henk/was_00002.jpg",
		"/henk/i am unique",
		"/henk/youKnowMe",
		"/henk/youKNOWMe_00001",
	}, efs2) {
		return
	}
}
