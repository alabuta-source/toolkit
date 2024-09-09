package cryptbuilder

import (
	"bytes"
	"encoding/gob"
	"strings"
)

// EncodeFromData use to encode your data before encrypt it.
// remember to register your data before anything.
// must be called on runtime.
//
// exemple:
//
//	func init() {
//	    gob.Register({your data pointer})
//	}
func EncodeFromData(data any) (string, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func DecodeFromBuffer(encryptedData string, decode any) error {
	reader := strings.NewReader(encryptedData)
	if err := gob.NewDecoder(reader).Decode(decode); err != nil {
		return err
	}
	return nil
}
