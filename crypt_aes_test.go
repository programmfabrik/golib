package golib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptAES(t *testing.T) {
	text := "henk"
	secret := "12345678901234567890123456789012"

	cText, err := EncryptAES(text, secret)
	if !assert.NoError(t, err) {
		return
	}

	pText, err := DecryptAES(cText, secret)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.Equal(t, text, pText) {
		return
	}
}
