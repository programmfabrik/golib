package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCut(t *testing.T) {
	s := []string{
		"hello this is a test",
		`multiline
        info`,
		"123456789012345",
	}

	e := []string{
		"hello this is a...",
		"multiline⏎inf...",
		"123456789012345",
	}

	assert.Equal(t, e, CutStrInArray(s, 15, "..."))
}
