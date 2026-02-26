package usecases

import (
	"context"
	"fmt"
	"time"
	"typeten/internal/domain"
)

// createSessionUserRepo defines the user repository methods needed for this use case.
type createSessionUserRepo interface {
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)
}

// createSessionTextRepo defines the text repository methods needed for this use case.
type createSessionTextRepo interface {
	GetTextInfo(ctx context.Context, id domain.TextID) (*domain.TextInfo, error)
}

// createSessionSessionRepo defines the session repository methods needed for this use case.
type createSessionSessionRepo interface {
	Create(ctx context.Context, session *domain.Session) error
}

// CreateSessionUseCase handles creating a new typing session.
type CreateSessionUseCase struct {
	sessionRepo createSessionSessionRepo
	textRepo    createSessionTextRepo
	userRepo    createSessionUserRepo
}

// NewCreateSessionUseCase creates a new CreateSessionUseCase.
func NewCreateSessionUseCase(sessionRepo createSessionSessionRepo, textRepo createSessionTextRepo, userRepo createSessionUserRepo) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		sessionRepo: sessionRepo,
		textRepo:    textRepo,
		userRepo:    userRepo,
	}
}

// CreateSessionInput represents the input for creating a session.
type CreateSessionInput struct {
	UserID domain.UserID
	TextID domain.TextID
}

// CreateSessionOutput represents the result of creating a session.
type CreateSessionOutput struct {
	Session *domain.Session
}

// Execute creates a new session after validating user and text exist.
func (uc *CreateSessionUseCase) Execute(ctx context.Context, input CreateSessionInput) (*CreateSessionOutput, error) {
	// Verify user exists
	_, err := uc.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Verify text exists
	_, err = uc.textRepo.GetTextInfo(ctx, input.TextID)
	if err != nil {
		return nil, fmt.Errorf("text not found: %w", err)
	}

	// Create session
	sessionID := domain.SessionID(fmt.Sprintf("session_%d", time.Now().UnixNano()))
	now := time.Now()

	session, err := domain.NewSession(sessionID, input.UserID, input.TextID, now)
	if err != nil {
		return nil, err
	}

	// Store session
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &CreateSessionOutput{Session: session}, nil
}
