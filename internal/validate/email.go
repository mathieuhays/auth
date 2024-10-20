package validate

import (
	"errors"
	"net/mail"
	"strings"
)

var ErrInvalidEmail = errors.New("invalid email")

func Email(email string) error {
	// no domain (or too many)
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ErrInvalidEmail
	}

	// missing tld
	if !strings.Contains(parts[1], ".") {
		return ErrInvalidEmail
	}

	// failsafe in case it catches something we don't
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}

	return nil
}
