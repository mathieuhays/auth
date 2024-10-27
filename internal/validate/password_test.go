package validate

import (
	"errors"
	"testing"
)

func TestPassword(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		err   error
	}{
		{"empty", "", ErrPasswordEmpty},
		{"false empty", "   ", ErrPasswordEmpty},
		{"too short", "test", ErrPasswordTooShort},
		{"too short padded", "  test  ", ErrPasswordTooShort},
		{"good pass", "test1234", nil},
		{"good pass padded", "  test1234  ", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Password(tc.input)

			if !errors.Is(err, tc.err) {
				t.Errorf("invalid output. expected: %v. got: %v", tc.err, err)
			}
		})
	}
}
