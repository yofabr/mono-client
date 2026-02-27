package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash converts plain text into a bcrypt hash (salt is handled by bcrypt).
func Hash(text string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare validates that text matches the stored bcrypt hash.
func Compare(text string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(text))
}

// KeyCacheUserID returns the canonical Redis key for a user auth session.
func KeyCacheUserID(userId string) string {
	return "auth:user:" + userId
}
