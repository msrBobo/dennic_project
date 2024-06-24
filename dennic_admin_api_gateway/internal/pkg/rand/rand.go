package rand

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateRandomNumbers() string {
	n := 5
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		code += strconv.Itoa(rand.Intn(9))
	}
	return code
}
