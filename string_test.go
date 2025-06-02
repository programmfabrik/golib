package golib

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/collate"
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

func TestSortStr2(t *testing.T) {
	lang := language.Make("de-DE")
	s1 := SortStr(lang, "000000001-234", collate.Numeric)
	s2 := SortStr(lang, "20201007-123", collate.Numeric)
	s3 := SortStr(lang, "00000000-123", collate.Numeric)

	s := []string{s2, s3, s1}
	sort.Strings(s)
	if !assert.Equal(t, []string{s1, s2, s3}, s) {
		return
	}

	s = []string{s1, s2}
	sort.Strings(s)
	if !assert.Equal(t, []string{s1, s2}, s) {
		return
	}
}

// 14f0049c14f0000114e7        0000002000200020002000200020
// 14f0000114e7049c14f0000114e700000020002000200020002000200020
