package validate

import (
	"errors"
	"strings"
)

var (
	ErrPasswordEmpty    = errors.New("password is empty")
	ErrPasswordTooShort = errors.New("password is too short")
)

func Password(value string) error {
	password := strings.Trim(value, " ")

	if password == "" {
		return ErrPasswordEmpty
	}

	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}
