package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func encrypt(s string) string {
	return "nil"
}

func decrypt(s string) string {
	return "nil"
}

func encryptSHA256(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))

	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	return mdStr
}
