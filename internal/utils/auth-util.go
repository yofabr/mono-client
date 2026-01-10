package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(text string) (string, error) {
	// GenerateFromPassword handles salt generation automatically
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	return string(bytes), err
}

func Compare(text string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(text))
}
