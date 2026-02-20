package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	now := time.Now()
	validSessionID := SessionID("session_123")
	validUserID := UserID("user_123")
	validTextID := TextID("text_123")

	tests := []struct {
		name      string
		id        SessionID
		userID    UserID
		textID    TextID
		now       time.Time
		wantErr   error
	}{
		{
			name:    "valid session",
			id:      validSessionID,
			userID:  validUserID,
			textID:  validTextID,
			now:     now,
			wantErr: nil,
		},
		{
			name:    "empty session ID",
			id:      "",
			userID:  validUserID,
			textID:  validTextID,
			now:     now,
			wantErr: ErrInvalidID,
		},
		{
			name:    "empty user ID",
			id:      validSessionID,
			userID:  "",
			textID:  validTextID,
			now:     now,
			wantErr: ErrInvalidID,
		},
		{
			name:    "empty text ID",
			id:      validSessionID,
			userID:  validUserID,
			textID:  "",
			now:     now,
			wantErr: ErrInvalidID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := NewSession(tt.id, tt.userID, tt.textID, tt.now)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if session == nil {
					t.Fatal("NewSession() returned nil")
				}
				if session.ID != tt.id {
					t.Errorf("NewSession() ID = %v, want %v", session.ID, tt.id)
				}
				if session.CurrentFragmentIdx != 0 {
					t.Errorf("NewSession() CurrentFragmentIdx = %v, want 0", session.CurrentFragmentIdx)
				}
				if session.CompletedLines != 0 {
					t.Errorf("NewSession() CompletedLines = %v, want 0", session.CompletedLines)
				}
				if session.IsCompleted {
					t.Error("NewSession() IsCompleted = true, want false")
				}
			}
		})
	}
}

func TestSession_RecordLineCompleted(t *testing.T) {
	now := time.Now()
	session, err := NewSession("session_1", "user_1", "text_1", now)
	if err != nil {
		t.Fatalf("NewSession() error = %v", err)
	}

	tests := []struct {
		name            string
		accuracyPercent float64
		wpm             float64
		setup           func(*Session)
		wantErr         error
		wantCompleted   int
		wantAccuracy    float64
		wantWPM         float64
	}{
		{
			name:            "first line completion",
			accuracyPercent: 95.5,
			wpm:             45.2,
			setup:           func(s *Session) {},
			wantErr:         nil,
			wantCompleted:   1,
			wantAccuracy:    95.5,
			wantWPM:         45.2,
		},
		{
			name:            "second line completion",
			accuracyPercent: 98.0,
			wpm:             50.0,
			setup: func(s *Session) {
				s.CompletedLines = 1
				s.TotalAccuracyPercent = 95.5
				s.AverageWPM = 45.2
			},
			wantErr:       nil,
			wantCompleted: 2,
			wantAccuracy:  (95.5*1 + 98.0) / 2,
			wantWPM:       (45.2*1 + 50.0) / 2,
		},
		{
			name:            "negative accuracy",
			accuracyPercent: -1.0,
			wpm:             45.0,
			setup:           func(s *Session) {},
			wantErr:         ErrInvalidSessionOp,
		},
		{
			name:            "accuracy over 100",
			accuracyPercent: 101.0,
			wpm:             45.0,
			setup:           func(s *Session) {},
			wantErr:         ErrInvalidSessionOp,
		},
		{
			name:            "negative WPM",
			accuracyPercent: 95.0,
			wpm:             -1.0,
			setup:           func(s *Session) {},
			wantErr:         ErrInvalidSessionOp,
		},
		{
			name:            "already completed session",
			accuracyPercent: 95.0,
			wpm:             45.0,
			setup: func(s *Session) {
				s.IsCompleted = true
			},
			wantErr: ErrInvalidSessionOp,
		},
		{
			name:            "zero accuracy",
			accuracyPercent: 0.0,
			wpm:             45.0,
			setup:           func(s *Session) {},
			wantErr:         nil,
			wantCompleted:   1,
			wantAccuracy:    0.0,
			wantWPM:         45.0,
		},
		{
			name:            "zero WPM",
			accuracyPercent: 95.0,
			wpm:             0.0,
			setup:           func(s *Session) {},
			wantErr:         nil,
			wantCompleted:   1,
			wantAccuracy:    95.0,
			wantWPM:         0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				ID:                 session.ID,
				UserID:             session.UserID,
				TextID:             session.TextID,
				CurrentFragmentIdx: 0,
				CurrentLineIdx:     0,
				CompletedLines:     0,
				TotalAccuracyPercent: 0,
				AverageWPM:         0,
				IsCompleted:        false,
				CreatedAt:          now,
				UpdatedAt:          now,
			}
			tt.setup(s)

			err := s.RecordLineCompleted(tt.accuracyPercent, tt.wpm, now.Add(time.Second))
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("RecordLineCompleted() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if s.CompletedLines != tt.wantCompleted {
					t.Errorf("RecordLineCompleted() CompletedLines = %v, want %v", s.CompletedLines, tt.wantCompleted)
				}
				if s.TotalAccuracyPercent != tt.wantAccuracy {
					t.Errorf("RecordLineCompleted() TotalAccuracyPercent = %v, want %v", s.TotalAccuracyPercent, tt.wantAccuracy)
				}
				if s.AverageWPM != tt.wantWPM {
					t.Errorf("RecordLineCompleted() AverageWPM = %v, want %v", s.AverageWPM, tt.wantWPM)
				}
			}
		})
	}

	// Test nil session
	var nilSession *Session
	err = nilSession.RecordLineCompleted(95.0, 45.0, now)
	if !errors.Is(err, ErrInvalidSessionOp) {
		t.Errorf("RecordLineCompleted() on nil session error = %v, wantErr %v", err, ErrInvalidSessionOp)
	}
}

func TestSession_MarkCompleted(t *testing.T) {
	now := time.Now()
	session, err := NewSession("session_1", "user_1", "text_1", now)
	if err != nil {
		t.Fatalf("NewSession() error = %v", err)
	}

	tests := []struct {
		name    string
		setup   func(*Session)
		wantErr error
	}{
		{
			name:    "mark as completed",
			setup:   func(s *Session) {},
			wantErr: nil,
		},
		{
			name: "already completed",
			setup: func(s *Session) {
				s.IsCompleted = true
			},
			wantErr: ErrInvalidSessionOp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				ID:                 session.ID,
				UserID:             session.UserID,
				TextID:             session.TextID,
				CurrentFragmentIdx: 0,
				CurrentLineIdx:     0,
				CompletedLines:     0,
				TotalAccuracyPercent: 0,
				AverageWPM:         0,
				IsCompleted:        false,
				CreatedAt:          now,
				UpdatedAt:          now,
			}
			tt.setup(s)

			err := s.MarkCompleted(now.Add(time.Second))
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkCompleted() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if !s.IsCompleted {
					t.Error("MarkCompleted() IsCompleted = false, want true")
				}
			}
		})
	}

	// Test nil session
	var nilSession *Session
	err = nilSession.MarkCompleted(now)
	if !errors.Is(err, ErrInvalidSessionOp) {
		t.Errorf("MarkCompleted() on nil session error = %v, wantErr %v", err, ErrInvalidSessionOp)
	}
}
