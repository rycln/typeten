package domain

import "time"

type UserID string

type User struct {
	ID        UserID
	Email     string
	Username  string
	Stats     UserStats
	Settings  UserSettings
	CreatedAt time.Time
}

type UserStats struct {
	AverageAccuracy float64
	AverageWPM      float64
}

type UserSettings struct {
	Theme string
}
