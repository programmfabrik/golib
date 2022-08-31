package golib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptAES encrypts text to cryptedText using the secretKey. The secretKey must
// be exactly 32 bytes long. cryptedText is base64 encoded.
func EncryptAES(text, secretKey string) (cryptedText string, err error) {

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher([]byte(secretKey))
	// if there are any errors, handle them
	if err != nil {
		return "", fmt.Errorf("Creating cipher failed: %w", err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		return "", fmt.Errorf("Creating GCM failed: %w", err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("Creating Nonce failed: %w", err)
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(text), nil)), nil
}

// DecryptAES takes a base64 encoded cryptedText and decrypt it to text using
// provided secretKey. The secretKey must be exactly 32 bytes long. cryptedText
// is base64 encoded.
func DecryptAES(cryptedText, secretKey string) (text string, err error) {
	c, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("Creating cipher failed: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", fmt.Errorf("Creating GCM failed: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(cryptedText) < nonceSize {
		return "", fmt.Errorf("Nonce wrong length: %w", err)
	}

	cryptedBytes, err := base64.StdEncoding.DecodeString(cryptedText)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	nonce, cryptedBytes := cryptedBytes[:nonceSize], cryptedBytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cryptedBytes, nil)
	if err != nil {
		return "", fmt.Errorf("Unable to decrypt: %w", err)
	}
	return string(plaintext), nil
}
