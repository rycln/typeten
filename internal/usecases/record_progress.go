package usecases

import (
	"context"
	"fmt"
	"time"
	"typeten/internal/domain"
	"typeten/internal/repository"
)

// RecordProgressUseCase handles recording typing progress for a session.
type RecordProgressUseCase struct {
	sessionRepo repository.SessionRepository
}

// NewRecordProgressUseCase creates a new RecordProgressUseCase.
func NewRecordProgressUseCase(sessionRepo repository.SessionRepository) *RecordProgressUseCase {
	return &RecordProgressUseCase{
		sessionRepo: sessionRepo,
	}
}

// RecordProgressInput represents the input for recording progress.
type RecordProgressInput struct {
	SessionID       string
	AccuracyPercent float64
	WPM             float64
}

// RecordProgressOutput represents the result of recording progress.
type RecordProgressOutput struct {
	Session *domain.Session
}

// Execute records a completed line and updates session statistics.
func (uc *RecordProgressUseCase) Execute(ctx context.Context, input RecordProgressInput) (*RecordProgressOutput, error) {
	// Get session
	session, err := uc.sessionRepo.GetByID(ctx, domain.SessionID(input.SessionID))
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	
	// Record line completion
	now := time.Now()
	if err := session.RecordLineCompleted(input.AccuracyPercent, input.WPM, now); err != nil {
		return nil, fmt.Errorf("failed to record progress: %w", err)
	}
	
	// Update session in repository
	if err := uc.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}
	
	return &RecordProgressOutput{Session: session}, nil
}
