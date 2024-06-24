package v1

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber(length int) string {
	rand.Seed(time.Now().UnixNano())
	const digits = "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}
