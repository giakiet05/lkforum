package util

import (
	"math/rand"
	"time"
)

func init() {
	// Seed the random number generator once when the package is initialized.
	// This is sufficient for non-cryptographic random strings.
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// GenerateRandomString generates a random string of a given length.
func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
