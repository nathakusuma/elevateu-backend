package uuidpkg

import (
	"sync"

	"github.com/google/uuid"
)

type IUUID interface {
	NewV7() (uuid.UUID, error)
}

type uuidStruct struct{}

var (
	uuidInstance IUUID
	once         sync.Once
)

func GetUUID() IUUID {
	once.Do(func() {
		uuidInstance = &uuidStruct{}
	})

	return uuidInstance
}

func (u *uuidStruct) NewV7() (uuid.UUID, error) {
	return uuid.NewV7()
}
