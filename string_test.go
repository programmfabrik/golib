package golib

import (
	"encoding/hex"
	"slices"
	"sort"
	"strings"
	"sync"
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
	s1 := SortStr(lang, "00000000-123", collate.Numeric)
	s2 := SortStr(lang, "000000001-234", collate.Numeric)
	s3 := SortStr(lang, "20201007-123", collate.Numeric)

	s := []string{s2, s3, s1}

	sort.Strings(s)
	if !assert.Equal(t, []string{s1, s2, s3}, s) {
		return
	}

	s = []string{s2, s1}
	sort.Strings(s)
	if !assert.Equal(t, []string{s1, s2}, s) {
		return
	}
}

// TestSortStrCached checks that the cached collator paths return the same
// keys as freshly built collators (including numeric sorting requested by the
// language tag instead of collate.Numeric), also when used concurrently.
func TestSortStrCached(t *testing.T) {
	inputs := []string{"café", "café", "00000000-123", "20201007-123", "2 zwei", "10 zehn", "Ähre"}

	uncached := func(lang language.Tag, s string, opts ...collate.Option) string {
		opts = append(opts, collate.IgnoreWidth, collate.IgnoreCase)
		cl := collate.New(lang, opts...)
		if cl.CompareString("2", "10") < 0 {
			s = zeroRun.ReplaceAllString(s, `$1₀$3`)
		}
		buf := new(collate.Buffer)
		return hex.EncodeToString(cl.KeyFromString(buf, s))
	}

	for _, lang := range []language.Tag{
		language.Make("de-DE"),
		language.Make("en"),
		language.Make("de-DE-u-kn-true"), // numeric by tag
	} {
		for _, s := range inputs {
			assert.Equal(t, uncached(lang, s), SortStr(lang, s), "plain %s %q", lang, s)
			assert.Equal(t, uncached(lang, s, collate.Numeric), SortStrNumeric(lang, s), "numeric %s %q", lang, s)
		}
	}

	lang := language.Make("de-DE")
	wg := sync.WaitGroup{}
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 500 {
				for _, s := range inputs {
					assert.Equal(t, uncached(lang, s, collate.Numeric), SortStrNumeric(lang, s))
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkSortStrNumeric(b *testing.B) {
	lang := language.Make("de-DE")
	for b.Loop() {
		SortStrNumeric(lang, "20201007-123 example title")
	}
}

func TestSortStr3(t *testing.T) {
	lang := language.Make("de-DE")

	sortStrings := func(ss []string, opts ...collate.Option) {
		slices.SortStableFunc(ss, func(a, b string) int {
			return strings.Compare(SortStr(lang, a, opts...), SortStr(lang, b, opts...))
		})
	}

	s := []string{"2 zwei", "1 eins", "11 elf", "100 hundert", "10 zehn", "9 yo", "0 null"}
	// cl := collate.New(language.Make("de-DE"), collate.IgnoreWidth, collate.IgnoreCase, collate.Numeric)
	// cl.SortStrings(s)
	sortStrings(s, collate.Numeric)
	if !assert.Equal(t, []string{"0 null", "1 eins", "2 zwei", "9 yo", "10 zehn", "11 elf", "100 hundert"}, s) {
		return
	}
	sortStrings(s)
	if !assert.Equal(t, []string{"0 null", "1 eins", "10 zehn", "100 hundert", "11 elf", "2 zwei", "9 yo"}, s) {
		return
	}

	s = []string{"20200901-123", "000000001-234", "00000001-123"}
	sortStrings(s, collate.Numeric)
	if !assert.Equal(t, []string{"00000001-123", "000000001-234", "20200901-123"}, s) {
		return
	}
}

// 14f0049c14f0000114e7        0000002000200020002000200020
// 14f0000114e7049c14f0000114e700000020002000200020002000200020
