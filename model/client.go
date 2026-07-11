package model

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	clientNameMinLength     = 3
	clientNameMaxLength     = 255
	clientEmailMaxLength    = 255
	clientPasswordMinLength = 8
	clientPasswordMaxBytes  = 72
)

var clientEmailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

type Client struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewClient(name, email, passwordHash string) (Client, error) {
	c := Client{
		ID:           uuid.New(),
		Name:         name,
		Email:        NormalizeClientEmail(email),
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	if err := validateClientName(c.Name); err != nil {
		return Client{}, err
	}
	if err := validateClientEmail(c.Email); err != nil {
		return Client{}, err
	}
	if err := validateClientPasswordHash(c.PasswordHash); err != nil {
		return Client{}, err
	}

	return c, nil
}

func NormalizeClientEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func validateClientName(name string) error {
	if name == "" {
		return ErrClientNameRequired
	}
	if utf8.RuneCountInString(name) < clientNameMinLength {
		return ErrClientNameTooShort
	}
	if utf8.RuneCountInString(name) > clientNameMaxLength {
		return ErrClientNameTooLong
	}
	return nil
}

func ValidateClientEmail(email string) error {
	if email == "" {
		return ErrClientEmailRequired
	}
	if !clientEmailRegex.MatchString(email) {
		return ErrClientEmailInvalid
	}
	if utf8.RuneCountInString(email) > clientEmailMaxLength {
		return ErrClientEmailTooLong
	}
	return nil
}

func validateClientEmail(email string) error {
	return ValidateClientEmail(email)
}

func ValidateClientPassword(password string) error {
	if password == "" {
		return ErrClientPasswordRequired
	}
	if utf8.RuneCountInString(password) < clientPasswordMinLength {
		return ErrClientPasswordTooShort
	}
	if len(password) > clientPasswordMaxBytes {
		return ErrClientPasswordTooLong
	}
	return nil
}

func validateClientPasswordHash(hash string) error {
	if hash == "" {
		return ErrClientPasswordHashRequired
	}
	return nil
}
