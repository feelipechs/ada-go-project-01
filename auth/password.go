package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordTooLong = errors.New("password exceeds maximum length")

func HashPassword(password string) (string, error) {
	if len(password) > 72 {
		return "", ErrPasswordTooLong
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password, passwordHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}
