package domain

type SessionID string

type Session struct {
	ID                   SessionID
	UserID               UserID
	TextID               TextID
	CurrentFragmentIdx   int
	CurrentLineIdx       int
	CompletedLines       int
	TotalAccuracyPercent float64
	AverageWPM           float64
	IsCompleted          bool
}
