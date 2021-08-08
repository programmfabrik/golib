package lib

import (
	"github.com/gofrs/uuid"
)

func NewUUID() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return uuid.String()
}
