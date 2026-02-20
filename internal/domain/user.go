// Package domain defines the core models and business rules for the typing
// application: users, texts (and their fragments), and typing sessions.
package domain

import (
	"regexp"
	"strings"
	"time"
)

// UserID identifies a user. It must be non-empty when constructing a User.
type UserID string

// User is the aggregate for an account. Email and Username must be valid and non-empty.
type User struct {
	ID        UserID
	Email     string
	Username  string
	CreatedAt time.Time
}

// Minimal email pattern: non-empty local part, @, non-empty domain with at least one dot.
var emailRx = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)

// NewUser creates a User after validating ID, email, and username.
// Returns ErrInvalidUser (wrapping ErrInvalidID) if any field is invalid.
func NewUser(id UserID, email, username string, createdAt time.Time) (*User, error) {
	if err := validateUserID(id); err != nil {
		return nil, err
	}
	email = strings.TrimSpace(email)
	username = strings.TrimSpace(username)
	if email == "" || !emailRx.MatchString(email) {
		return nil, ErrInvalidUser
	}
	if username == "" {
		return nil, ErrInvalidUser
	}
	return &User{
		ID:        id,
		Email:     email,
		Username:  username,
		CreatedAt: createdAt,
	}, nil
}

func validateUserID(id UserID) error {
	if strings.TrimSpace(string(id)) == "" {
		return ErrInvalidID
	}
	return nil
}
