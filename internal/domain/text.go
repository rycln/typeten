package domain

import "time"

type (
	TextID         string
	TextFragmentID string
)

type TextInfo struct {
	ID            TextID
	UserID        UserID
	Title         string
	TotalLines    int
	FragmentSize  int
	FragmentCount int
	CreatedAt     time.Time
}

type TextFragment struct {
	ID          TextFragmentID
	TextID      TextID
	FragmentIdx int
	Lines       []string
}
