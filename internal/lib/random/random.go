package random

import (
	"math/rand"
)

var alphabet = []rune("abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"1234567890")

func NewRandomString(length int) string {
	response := make([]rune, length)

	for i := range response {
		response[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(response)
}
