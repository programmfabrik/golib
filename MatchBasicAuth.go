package lib

import (
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// MatchBasicAuth checks if a request contains credentials matched in basicAuth map.
// "auth" is the http request header Authorization
// It returns true if matched, false otherwise. If basicAuth is empty, this check returns true.
func MatchBasicAuth(auth string, basicAuth map[string]string) bool {

	s := strings.SplitN(auth, " ", 2)
	if len(s) != 2 {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}

	for user, pass := range basicAuth {
		if pair[0] != user {
			continue
		}
		// try matching without bcrypt
		if pair[1] == pass {
			return true
		}
		// try bcrypt
		err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(pair[1]))
		// if err != nil {
		// 	println(err.Error())
		// }
		if err != nil {
			return false
		}
		return true
	}

	return false
}
