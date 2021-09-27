package golib

import (
	"crypto/rand"
	mrand "math/rand"
)

// RandStr creates a random string containg of n chars
// seed with rand.Seed(time.Now().UnixNano())
func RandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

const otpChars = "0123456789"

func GenerateNumberOTP(length int) (string, error) {
	return GenerateOTP(length, otpChars)
}

// GenerateOTP creates a random string from the
// given chars
func GenerateOTP(length int, otpChars string) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
