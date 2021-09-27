package golib

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
)

// Max returns the larger of x or y.
func IntMax(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func IntMin(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Pln(s string, args ...interface{}) {
	println(fmt.Sprintf(s, args...))
}

type StringMap map[string]interface{}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

// GetInt returns int64 from a string, default to the given default if 0 or empty
func GetInt(s string, d int64) int64 {
	if s == "" || s == "0" {
		return d
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	if i == 0 {
		return d
	}
	return i
}

// GetBool returns bool from a string
func GetBool(s string) bool {
	if s == "true" || s == "1" {
		return true
	}
	return false
}

func ParseInt64(input string) (int64, error) {
	flt, _, err := big.ParseFloat(input, 10, 0, big.ToNearestEven)
	if err != nil {
		return 0, err
	}
	var i = new(big.Int)
	i, _ = flt.Int(i)
	return i.Int64(), nil
}

type Date struct {
	Value *string `json:"value"`
}
