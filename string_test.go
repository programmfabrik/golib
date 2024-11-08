package golib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
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

func TestSortStr(t *testing.T) {
	lang := language.Make("de-DE")
	s_nfd := SortStr(lang, "cafe\u0301") // UTF-8: \x63\x61\x66\x65\xcc\x81)
	s_nfc := SortStr(lang, "caf\u00e9")  // UTF-8: \x63\x61\x66\xc3\xa9
	if !assert.Equal(t, s_nfd, s_nfc, "NFD / NFC not equal") {
		return
	}

	s_nfd = SortStr(lang, "cr\u00e8me br\u00fbl\u00e9e")    // NFC crème brûlée
	s_nfc = SortStr(lang, "cre\u0300me bru\u0302le\u0301e") // NFD crème brûlée
	if !assert.Equal(t, s_nfd, s_nfc, "NFD / NFC not equal") {
		return
	}
}
