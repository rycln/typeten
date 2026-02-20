package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	now := time.Now()
	validID := UserID("user_123")
	validEmail := "test@example.com"
	validUsername := "testuser"

	tests := []struct {
		name    string
		id      UserID
		email   string
		username string
		createdAt time.Time
		wantErr error
	}{
		{
			name:     "valid user",
			id:       validID,
			email:    validEmail,
			username: validUsername,
			createdAt: now,
			wantErr:  nil,
		},
		{
			name:     "empty ID",
			id:       "",
			email:    validEmail,
			username: validUsername,
			createdAt: now,
			wantErr:  ErrInvalidID,
		},
		{
			name:     "whitespace ID",
			id:       "   ",
			email:    validEmail,
			username: validUsername,
			createdAt: now,
			wantErr:  ErrInvalidID,
		},
		{
			name:     "empty email",
			id:       validID,
			email:    "",
			username: validUsername,
			createdAt: now,
			wantErr:  ErrInvalidUser,
		},
		{
			name:     "invalid email format",
			id:       validID,
			email:    "notanemail",
			username: validUsername,
			createdAt: now,
			wantErr:  ErrInvalidUser,
		},
		{
			name:     "email without domain",
			id:       validID,
			email:    "test@",
			username: validUsername,
			createdAt: now,
			wantErr:  ErrInvalidUser,
		},
		{
			name:     "empty username",
			id:       validID,
			email:    validEmail,
			username: "",
			createdAt: now,
			wantErr:  ErrInvalidUser,
		},
		{
			name:     "whitespace username",
			id:       validID,
			email:    validEmail,
			username: "   ",
			createdAt: now,
			wantErr:  ErrInvalidUser,
		},
		{
			name:     "email with whitespace trimmed",
			id:       validID,
			email:    "  test@example.com  ",
			username: validUsername,
			createdAt: now,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.id, tt.email, tt.username, tt.createdAt)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if user == nil {
					t.Fatal("NewUser() returned nil user")
				}
				if user.ID != tt.id {
					t.Errorf("NewUser() ID = %v, want %v", user.ID, tt.id)
				}
				if user.Email != strings.TrimSpace(tt.email) {
					t.Errorf("NewUser() Email = %v, want %v", user.Email, strings.TrimSpace(tt.email))
				}
				if user.Username != strings.TrimSpace(tt.username) {
					t.Errorf("NewUser() Username = %v, want %v", user.Username, strings.TrimSpace(tt.username))
				}
			}
		})
	}
}
