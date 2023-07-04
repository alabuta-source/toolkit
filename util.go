package toolkit

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"path"
	"runtime"
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

func generateUUID() string {
	return uuid.New().String()
}

func getRootDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to read current file")
	}
	pwd := path.Dir(filename)
	return path.Join(pwd, "..", ".."), nil
}

func formatErr(msg string, args ...interface{}) string {
	return fmt.Sprintf(msg, args...)
}

func cutSpaces(value string) string {
	return strings.Replace(value, " ", "", -1)
}

func removeBucketName(path, bucket string) string {
	return strings.Replace(path, fmt.Sprintf("/%s/", bucket), "", -1)
}
