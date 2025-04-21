package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GenerateRef(input string, value string, truncate string) string {
	key := []byte(input)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(value))
	hashed := hex.EncodeToString(h.Sum(nil))
	if truncate == "hard" {
		return hashed[12:44]
	} else if truncate == "soft" {
		return hashed[0:16]
	}
	return hashed
}

func BuildReferenceString(parts ...string) string {
	return strings.Join(parts, "")
}