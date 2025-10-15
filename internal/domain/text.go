package domain

import "time"

type TextID string

type TextInfo struct {
	ID         TextID
	UserID     UserID
	Title      string
	TotalLines int
	CreatedAt  time.Time
}

type TextFragmentID string

type TextFragment struct {
	ID            TextFragmentID
	TextID        TextID
	FragmentIndex int
	Content       string
}
