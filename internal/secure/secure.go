package secure

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
)

func generateRandomSalt(saltSize int) []byte {
	salt := make([]byte, saltSize)
	_, _ = rand.Read(salt)
	return salt
}

func GenerateHashWithSalt(v string, salt []byte) string {
	bytes := []byte(v)
	hasher := sha512.New()
	bytes = append(bytes, salt...)
	_, _ = hasher.Write(bytes)
	hash := hasher.Sum(nil)
	return base64.URLEncoding.EncodeToString(hash)
}

func GenerateHash(v string) string {
	const saltSize = 64
	return GenerateHashWithSalt(v, generateRandomSalt(saltSize))
}
