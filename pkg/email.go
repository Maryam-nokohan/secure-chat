package pkg

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return errors.New("email is required")
	}

	if len(email) > 254 {
		return errors.New("email address is too long")
	}

	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return errors.New("invalid email address")
	}

	return nil
}