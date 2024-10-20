package validate

import (
	"errors"
	"testing"
)

func TestEmail(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		err   error
	}{
		{"empty", "", ErrInvalidEmail},
		{"random string", "invalid", ErrInvalidEmail},
		{"missing tld", "test@domain", ErrInvalidEmail},
		{"valid email", "test@example.com", nil},
		{"valid email with alias", "test+123@example.com", nil},
		{"valid email with dot in name", "test.account@example.com", nil},
		{"valid email with subdomain", "account@test.example.com", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Email(tc.input)

			if !errors.Is(err, tc.err) {
				t.Errorf("invalid output. expected: %v. got: %v", tc.err, err)
			}
		})
	}
}
