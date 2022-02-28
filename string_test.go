package golib

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

func TestPadStr(t *testing.T) {
	assert.Equal(t, "Ähre   ", PadStr("Ähre", 7))
	assert.Equal(t, "Ähr", PadStr("Ähre", 3))
	assert.Equal(t, "Ähre", PadStr("Ähre", 4))
}

func TestSplitChunks(t *testing.T) {
	assert.Equal(t, []string{"Ä", "Ö", "ß ", "Üa", "bde", "fg"}, StringByteChunks("ÄÖß Üabdefg", 3))
}
