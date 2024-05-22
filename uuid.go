package golib

import (
	"github.com/gofrs/uuid"
)

// NewUUID returns a randomly generated UUID
func NewUUID() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return uuid.String()
}

// IsUUID returns true if id is the same format as NewUUID provides
func IsUUID(id string) bool {
	uid, err := uuid.FromString(id)
	if err != nil {
		return false
	}
	if uid.Version() != uuid.V4 {
		return false
	}
	return true
}
