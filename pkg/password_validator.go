package pkg

import (
	"errors"
	"unicode"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New(
			"password must be at least 8 characters long",
		)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, ch := range password {

		switch {
		case unicode.IsUpper(ch):
			hasUpper = true

		case unicode.IsLower(ch):
			hasLower = true

		case unicode.IsDigit(ch):
			hasNumber = true

		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New(
			"password must contain at least one uppercase letter",
		)
	}

	if !hasLower {
		return errors.New(
			"password must contain at least one lowercase letter",
		)
	}

	if !hasNumber {
		return errors.New(
			"password must contain at least one number",
		)
	}

	if !hasSpecial {
		return errors.New(
			"password must contain at least one special character",
		)
	}

	return nil
}