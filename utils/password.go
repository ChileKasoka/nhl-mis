package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPass(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("err hashing password: %w", err)
	}
	return string(hash), nil
}

func ComparePass(password, hash) error {
	compare := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	return compare
}