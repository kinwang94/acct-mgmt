package utils

import (
	"math/rand"
	"strings"
	"time"
)

const (
	digits    = "0123456789"
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randomString(min, max int, charSets ...string) string {
	length := r.Intn(max-min+1) + min
	allChars := strings.Join(charSets, "")

	var result []rune
	if min >= len(charSets) {
		for _, set := range charSets {
			index := r.Intn(len(set))
			result = append(result, rune(set[index]))
		}
		length -= len(result)
	}

	for i := 0; i < length; i++ {
		index := r.Intn(len(allChars))
		result = append(result, rune(allChars[index]))
	}

	r.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return string(result)
}

// RandomStringLetter returns a random length string of letters (A-Z, a-z) between |min| and |max|.
func RandomStringLetter(min, max int) string {
	return randomString(min, max, lowercase, uppercase)
}

// RandomStringAlphanumeric returns a random length string of letters (A-Z, a-z) and digits (0-9) between |min| and |max|.
func RandomStringAlphanumeric(min, max int) string {
	return randomString(min, max, lowercase, uppercase, digits)
}
