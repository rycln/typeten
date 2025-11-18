package domain

import "time"

type UserID string

type User struct {
	ID        UserID
	Email     string
	Username  string
	CreatedAt time.Time
}
