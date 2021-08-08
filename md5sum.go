package lib

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

// MD5 JSONMarshals i and return the MD5 of the marshaled bytes
func MD5(i interface{}) string {
	bytes, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	hasher := md5.New()
	hasher.Write(bytes)
	return hex.EncodeToString(hasher.Sum(nil))
}
