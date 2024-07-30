package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomEmail() string {
	usernameLength := rand.Intn(5) + 5
	domainLength := rand.Intn(5) + 3
	tldLength := rand.Intn(3) + 2

	username := RandomString(usernameLength)
	domain := RandomString(domainLength)
	tld := RandomString(tldLength)

	return fmt.Sprintf("%s@%s.%s", username, domain, tld)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(10, 1000)
}

func RandomTransactionAmount() int64 {
	return RandomInt(-10, 10)
}

func RandomCurrency() string {
	currencies := []string{"INR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
