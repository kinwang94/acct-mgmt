package utils

import (
	"acct-mgmt/errors"
	"unicode"
)

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.ErrUsernameLengthTooShort
	}
	if len(username) > 32 {
		return errors.ErrUsernameLengthTooLong
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.ErrPasswordLengthTooShort
	}
	if len(password) > 32 {
		return errors.ErrPasswordLengthTooLong
	}

	hasUpper := false
	hasLower := false
	hasNumber := false

	for _, c := range password {
		if unicode.IsUpper(c) {
			hasUpper = true
		} else if unicode.IsLower(c) {
			hasLower = true
		} else if unicode.IsDigit(c) {
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return errors.ErrPasswordInvalidFormat
	}

	return nil
}
