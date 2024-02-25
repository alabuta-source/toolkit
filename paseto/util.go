package paseto

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func randomOwner() string {
	return randomString(6)
}

func formatErr(msg string, args ...interface{}) string {
	return fmt.Sprintf(msg, args...)
}

func setupOptions(option ...Option) *payloadOption {
	var options payloadOption
	for _, opt := range option {
		opt(&options)
	}
	return &options
}
