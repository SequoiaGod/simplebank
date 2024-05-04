package util

import (
	"fmt"
	"math/rand"
	"strings"
)

func RandomString(size int) string {
	const alphabet string = "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder
	for i := 0; i < size; i++ {
		sb.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}
	return sb.String()
}

func randomInt(min, max int64) int64 {
	return rand.Int63n(int64(max-min) + min)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return randomInt(10000, 100000)
}

func RandomCurrency() string {
	var currency = []string{"EUR", "USD", "GBP", "JPY", "CHY"}

	return currency[rand.Intn(len(currency))]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomString(6))
}
