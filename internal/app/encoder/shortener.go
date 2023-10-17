package encoder

import (
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Code() string {
	randomstring := make([]byte, 8)

	for i := 0; i < 8; i++ {
		index := rand.Intn(len(letters) - 1)

		randomstring[i] = letters[index]
	}
	code := string(randomstring)

	return code
}
