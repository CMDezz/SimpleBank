package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcefghijklmnpqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // 0 -> max - min
}

func RandomString(n int) string {
	var s strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		s.WriteByte(c)
	}
	return s.String()
}

func RandomOwner() string {
	return RandomString(6)
}
func RandomMoney() int64 {
	return RandomInt(200, 400)
}
func RandomCurrency() string {
	listCurr := []string{EUR, USD, VND}
	return listCurr[rand.Intn(len(listCurr))]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
