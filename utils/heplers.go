package utils

import (
	"math/big"
	"crypto/rand"
	"log"
)

func RandomFourDigits() int64 {
	max := big.NewInt(9999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
	}
	return n.Int64()
}

