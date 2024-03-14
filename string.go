package golib

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/cases"
)

// PushOntoStringArray adds str(s) to arr if missing
func PushOntoStringArray(arr []string, strs ...string) (arr1 []string) {
	arr1 = arr
	for _, str := range strs {
		if StrInArray(str, arr1) {
			continue
		}
		arr1 = append(arr1, str)
	}
	return arr1
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// CutStrInArray cuts each string inside the given array to the
// given length and adds the suffix is the string was cut.
// It also replaces newlines with ⏎ in the string
func CutStrInArray(arr []string, l int, suffix string) (arr2 []string) {
	arr2 = []string{}
	for _, s := range arr {
		s2 := []string{}
		for _, p := range strings.Split(s, "\n") {
			s2 = append(s2, strings.Trim(p, " "))
		}

		s = strings.Join(s2, "⏎")
		arr2 = append(arr2, CutStr(s, l, suffix))
	}
	return arr2
}

// CutStr cuts s if is longer than "len"
// When a string is cut, the suffix is added
func CutStr(s string, l int, suffix string) string {
	if len(s) > l {
		return s[0:l] + suffix
	}
	return s
}

// CutRunes works like CutStr but counts runes not bytes
func CutRunes(s string, l int, suffix string) string {
	r := []rune(s)
	if len(r) > l {
		return string(r[0:l]) + suffix
	}
	return s
}

// PadStr returns string s filled to a length of padWidth. If s is longer than pw
// string will be cut to the length. If padWidth <= 0 s will be returned unchanged.
// The counting of characters is rune based. So "Ä" counts as one.
func PadStr(s string, padWidth int) string {
	if padWidth <= 0 {
		return s
	}
	l := utf8.RuneCountInString(s)
	switch {
	case l == padWidth:
		return s
	case l < padWidth:
		return s + strings.Repeat(" ", padWidth-l)
	default:
		return string([]rune(s)[0:padWidth])
	}
}

func StrInArray(str string, arr []string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}

// ArrayContainsStrs returns true if arr contains any of the passed
// strings
func ArrayContainsStrs(arr []string, strs ...string) bool {
	for _, str := range strs {
		if StrInArray(str, arr) {
			return true
		}
	}
	return false
}

func ToString(i interface{}) string {
	if i == nil {
		return ""
	}
	switch v := i.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case *int64:
		return strconv.FormatInt(*v, 10)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", i)
	}
}

// ReplaceEndless replaces old to new in s as long
// as the string is getting shorter.
// This function can be used to remove double // from an URL path.
func ReplaceEndless(s, old, new string) string {
	for {
		ns := strings.ReplaceAll(s, old, new)
		if len(ns) >= len(s) {
			return s
		}
		s = ns
	}
}

// ToAnySlice converts a slice of []T to []any
func ToAnySlice[T any](list []T) []any {
	in := make([]any, len(list))
	for idx, s0 := range list {
		in[idx] = s0
	}
	return in
}

// AnyToStrSlice converts a slice of string to slice of interface
func AnyToStrSlice[T any](list []T) []string {
	sn := make([]string, len(list))
	for idx, i0 := range list {
		sn[idx] = fmt.Sprintf("%v", i0)
	}
	return sn
}

func FoldStr(s string) string {
	c := cases.Fold()
	return c.String(s)
}

// Split string s into byte chunks of a max size of chunkSize Each string
// returned has a maximum length of byteChunkSize The split is UTF-8 safe.
func StringByteChunks(s string, byteChunkSize int) (chunks []string) {
	if len(s) == 0 {
		return nil
	}
	if byteChunkSize >= len(s) {
		return []string{s}
	}
	chunks = []string{}
	chars := []byte{}
	chunkOffset := 0
	for idx, charRune := range s {
		char := string(charRune)
		if idx-chunkOffset+len(char) > byteChunkSize {
			// last char caused an overflow, add
			// to the previous last
			chunks = append(chunks, string(chars))
			chars = []byte{}
			chunkOffset = idx
		}
		chars = append(chars, []byte(char)...)
	}
	chunks = append(chunks, string(chars))
	return chunks
}

// ToValidUTF8 checks and converts a string to valid UTF-8, replacing invalid characters.
// It iterates over the string, and for each rune that is identified as invalid (RuneError),
// it replaces it with the specified 'replacement' rune.
func ToValidUTF8(s string, replacement rune) string {
	if utf8.ValidString(s) {
		return s // The string is already valid UTF-8.
	}

	// Build a new string with invalid runes replaced.
	var builder strings.Builder
	for _, r := range s {
		if r == utf8.RuneError {
			builder.WriteRune(replacement)
		} else {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

// DebugValues takes a slice of T and returns and ordered list each item
// rendered in a comma separated list. If the list is longer than length bytes,
// the rest of the slice is omitted and the output ends in ... With length <= 0,
// the whole slice is rendered
func DebugValues[T cmp.Ordered](list []T, length int) string {
	slices.Sort(list)
	sb := strings.Builder{}
	sb.WriteRune('[')
	for idx, item := range list {
		s := fmt.Sprintf("%v", item)
		if idx > 0 {
			sb.WriteRune(',')
		}
		if length > 0 && len(s)+sb.Len() > length {
			sb.WriteString("...")
			break
		}
		sb.WriteString(s)
	}
	sb.WriteString(fmt.Sprintf("] %d", len(list)))
	return sb.String()
}
