package utils

import (
	"math/rand"
	"time"
)

func createSeed(now time.Time) *rand.Rand {
	return rand.New(
		rand.NewSource(
			now.UnixNano(),
		),
	)
}

func stringWithCharset(now time.Time, charset string, length int) string {
	seed := createSeed(now)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}

	return string(b)
}

func GenerateOTP(now time.Time) string {
	charset := "0123456789"
	return stringWithCharset(now, charset, 4)
}
