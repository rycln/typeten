package domain

import (
	"strings"
	"time"
)

// SessionID identifies a typing session.
type SessionID string

// Session represents one user typing one text. CurrentFragmentIdx and CurrentLineIdx
// are the current typing position; CompletedLines is the number of lines fully
// completed. TotalAccuracyPercent and AverageWPM are running session-wide stats.
// Use RecordLineCompleted to update progress and MarkCompleted when the session ends.
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
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// NewSession creates a Session with validated IDs and non-negative indices/stats.
// Accuracy must be in [0, 100]; wpm must be >= 0. Returns ErrInvalidSession on invalid input.
func NewSession(id SessionID, userID UserID, textID TextID, now time.Time) (*Session, error) {
	if err := validateSessionID(id); err != nil {
		return nil, err
	}
	if err := validateUserID(userID); err != nil {
		return nil, err
	}
	if err := validateTextID(textID); err != nil {
		return nil, err
	}
	return &Session{
		ID:                   id,
		UserID:               userID,
		TextID:               textID,
		CurrentFragmentIdx:   0,
		CurrentLineIdx:       0,
		CompletedLines:       0,
		TotalAccuracyPercent: 0,
		AverageWPM:           0,
		IsCompleted:          false,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

// RecordLineCompleted updates CompletedLines and running averages for accuracy and WPM.
// accuracyPercent must be in [0, 100]; wpm must be >= 0. UpdatedAt is set to now.
// Returns ErrInvalidSessionOp if the session is already completed or values are invalid.
func (s *Session) RecordLineCompleted(accuracyPercent, wpm float64, now time.Time) error {
	if s == nil {
		return ErrInvalidSessionOp
	}
	if s.IsCompleted {
		return ErrInvalidSessionOp
	}
	if accuracyPercent < 0 || accuracyPercent > 100 || wpm < 0 {
		return ErrInvalidSessionOp
	}
	n := float64(s.CompletedLines)
	s.TotalAccuracyPercent = (s.TotalAccuracyPercent*n + accuracyPercent) / (n + 1)
	s.AverageWPM = (s.AverageWPM*n + wpm) / (n + 1)
	s.CompletedLines++
	s.UpdatedAt = now
	return nil
}

// MarkCompleted marks the session as completed and sets UpdatedAt to now.
// Returns ErrInvalidSessionOp if the session is nil or already completed.
func (s *Session) MarkCompleted(now time.Time) error {
	if s == nil {
		return ErrInvalidSessionOp
	}
	if s.IsCompleted {
		return ErrInvalidSessionOp
	}
	s.IsCompleted = true
	s.UpdatedAt = now
	return nil
}

func validateSessionID(id SessionID) error {
	if strings.TrimSpace(string(id)) == "" {
		return ErrInvalidID
	}
	return nil
}
