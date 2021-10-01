package golib

import (
	"fmt"
	"strconv"
	"strings"
)

func PushOntoStringArray(arr []string, str string) []string {
	for _, s := range arr {
		if s == str {
			return arr
		}
	}
	return append(arr, str)
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

// CurStr cuts s if is longer than "len"
// When a string is cut, the suffix is added
func CutStr(s string, l int, suffix string) string {
	if len(s) > l {
		return s[0:l] + suffix
	}
	return s
}

// PadStr returns string s filled to a length of padWidth. If s is longer than pw
// string will be cut to the length. If padWidth <= 0 s will be returned unchanged
func PadStr(s string, padWidth int) string {
	if padWidth <= 0 {
		return s
	}
	if len(s) < padWidth {
		return s + strings.Repeat(" ", padWidth-len(s))
	}
	return s[0:padWidth]
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

// StrSliceToInterfaceSlice converts a slice of string to slice of interface
func StrSliceToInterfaceSlice(s []string) []interface{} {
	in := make([]interface{}, len(s))
	for idx, s0 := range s {
		in[idx] = s0
	}
	return in
}

// InterfaceSliceToStrSlice converts a slice of string to slice of interface
func InterfaceSliceToStrSlice(i []interface{}) []string {
	sn := make([]string, len(i))
	for idx, i0 := range i {
		sn[idx] = fmt.Sprintf("%v", i0)
	}
	return sn
}
