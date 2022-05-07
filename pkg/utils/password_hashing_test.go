package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Run("TestHashPassword", func(t *testing.T) {
		hash := HashPassword("test.hash@password.com", "password")
		t.Logf("Hashed: %s", hash)
	})
}
