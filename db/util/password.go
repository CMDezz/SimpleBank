package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("err when hase password: %w", err)
	}
	return string(hashed_password), nil
}

func CheckPassword(password string, hased_password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hased_password), []byte(password))
}
