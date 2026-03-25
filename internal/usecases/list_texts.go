package usecases

import (
	"context"
	"fmt"
	"typeten/internal/domain"
)

// ListTextsUserRepo defines operations needed to validate a user exists.
type ListTextsUserRepo interface {
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)
}

// ListTextsTextRepo defines operations needed to list texts for a user.
type ListTextsTextRepo interface {
	ListByUserID(ctx context.Context, userID domain.UserID) ([]*domain.TextInfo, error)
}

// ListTextsUseCase handles listing texts for a user.
type ListTextsUseCase struct {
	textRepo ListTextsTextRepo
	userRepo ListTextsUserRepo
}

// NewListTextsUseCase creates a new ListTextsUseCase.
func NewListTextsUseCase(textRepo ListTextsTextRepo, userRepo ListTextsUserRepo) *ListTextsUseCase {
	return &ListTextsUseCase{
		textRepo: textRepo,
		userRepo: userRepo,
	}
}

// ListTextsInput represents the input for listing texts.
type ListTextsInput struct {
	UserID domain.UserID
}

// ListTextsOutput represents the result of listing texts.
type ListTextsOutput struct {
	Texts []*domain.TextInfo
}

// Execute lists all texts for a user.
func (uc *ListTextsUseCase) Execute(ctx context.Context, input ListTextsInput) (*ListTextsOutput, error) {
	// Verify user exists
	_, err := uc.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	texts, err := uc.textRepo.ListByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to list texts: %w", err)
	}
	
	return &ListTextsOutput{Texts: texts}, nil
}
