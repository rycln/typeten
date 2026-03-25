package usecases

import (
	"context"
	"fmt"
	"typeten/internal/domain"
)

// GetSessionSessionRepo defines operations needed to fetch a session.
type GetSessionSessionRepo interface {
	GetByID(ctx context.Context, id domain.SessionID) (*domain.Session, error)
}

// GetSessionUseCase handles retrieving a session by ID.
type GetSessionUseCase struct {
	sessionRepo GetSessionSessionRepo
}

// NewGetSessionUseCase creates a new GetSessionUseCase.
func NewGetSessionUseCase(sessionRepo GetSessionSessionRepo) *GetSessionUseCase {
	return &GetSessionUseCase{
		sessionRepo: sessionRepo,
	}
}

// GetSessionInput represents the input for getting a session.
type GetSessionInput struct {
	SessionID domain.SessionID
}

// GetSessionOutput represents the result of getting a session.
type GetSessionOutput struct {
	Session *domain.Session
}

// Execute retrieves a session by ID.
func (uc *GetSessionUseCase) Execute(ctx context.Context, input GetSessionInput) (*GetSessionOutput, error) {
	session, err := uc.sessionRepo.GetByID(ctx, input.SessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	
	return &GetSessionOutput{Session: session}, nil
}
