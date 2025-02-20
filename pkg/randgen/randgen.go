package randgen

import (
	"crypto/rand"
	"math"
	"math/big"
	"sync"
)

type IRandGen interface {
	RandomNumber(digits int) (int, error)
	RandomString(length int) (string, error)
}

type randGen struct{}

var (
	randGenInstance IRandGen
	once            sync.Once
)

func GetRandGen() IRandGen {
	once.Do(func() {
		randGenInstance = &randGen{}
	})

	return randGenInstance
}

func (r *randGen) RandomNumber(digits int) (int, error) {
	low := big.NewInt(int64(math.Pow10(digits - 1)))
	high := big.NewInt(int64(math.Pow10(digits) - 1))

	diff := new(big.Int).Sub(high, low)
	diff.Add(diff, big.NewInt(1)) // Add 1 to make it inclusive

	randomNum, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return 0, err
	}

	return int(randomNum.Add(randomNum, low).Int64()), nil
}

func (r *randGen) RandomString(length int) (string, error) {
	const alphaNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	alphaNumRunes := []rune(alphaNum)
	randomRune := make([]rune, length)

	for i := range length {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphaNumRunes))))
		if err != nil {
			return "", err
		}
		randomRune[i] = alphaNumRunes[randomIndex.Int64()]
	}

	return string(randomRune), nil
}
