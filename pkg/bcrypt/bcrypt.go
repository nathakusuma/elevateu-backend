package bcrypt

import (
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type IBcrypt interface {
	Hash(plain string) (string, error)
	Compare(password, hashed string) bool
}

type bcryptStruct struct{}

var (
	bcryptInstance IBcrypt
	once           sync.Once
)

func GetBcrypt() IBcrypt {
	once.Do(func() {
		bcryptInstance = &bcryptStruct{}
	})

	return bcryptInstance
}

func (b *bcryptStruct) Hash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)

	return string(bytes), err
}

func (b *bcryptStruct) Compare(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))

	return err == nil
}
