package utils

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"github.com/Xacnio/aptms-backend/pkg/configs"
)

func reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, _rune := range s {
		n--
		runes[n] = _rune
	}
	return string(runes[n:])
}

func HashPassword(email, password string) string {
	saltKey1 := configs.Get("SC_KEY1")
	saltKey2 := configs.Get("SC_KEY2")
	emailHash := NewSHA256([]byte(reverse(email) + saltKey1 + email + saltKey2))
	passwordHash := NewSHA256([]byte(reverse(password) + saltKey2 + password + saltKey1))
	hashedPassword := NewSHA256([]byte(saltKey1 + string(emailHash) + saltKey2 + string(passwordHash)))
	base64Hashed := b64.URLEncoding.EncodeToString(hashedPassword)
	return base64Hashed
}

func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
