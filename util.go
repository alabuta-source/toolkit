package toolkit

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func GenerateUUID() string {
	char1 := randomString(4)
	char2 := randomString(4)
	return fmt.Sprintf("%s-%s", char1, char2)
}

func IsValidPassword(pwd string) bool {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	for _, char := range pwd {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}
