package usecases

import (
	"context"
	"fmt"
	"time"
	"typeten/internal/domain"
)

// createTextUserRepo defines the user repository methods needed for this use case.
type createTextUserRepo interface {
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)
}

// createTextTextRepo defines the text repository methods needed for this use case.
type createTextTextRepo interface {
	CreateTextInfo(ctx context.Context, info *domain.TextInfo) error
	CreateFragment(ctx context.Context, fragment *domain.TextFragment) error
}

// CreateTextUseCase handles uploading and processing a new text.
type CreateTextUseCase struct {
	textRepo      createTextTextRepo
	userRepo      createTextUserRepo
	textProcessor *TextProcessor
}

// NewCreateTextUseCase creates a new CreateTextUseCase.
func NewCreateTextUseCase(textRepo createTextTextRepo, userRepo createTextUserRepo, fragmentSize int) *CreateTextUseCase {
	return &CreateTextUseCase{
		textRepo:      textRepo,
		userRepo:      userRepo,
		textProcessor: NewTextProcessor(fragmentSize),
	}
}

// CreateTextInput represents the input for creating a text.
type CreateTextInput struct {
	UserID  domain.UserID
	Title   string
	Content string
}

// CreateTextOutput represents the result of creating a text.
type CreateTextOutput struct {
	TextInfo *domain.TextInfo
}

// Execute creates a new text by processing the content and storing it.
func (uc *CreateTextUseCase) Execute(ctx context.Context, input CreateTextInput) (*CreateTextOutput, error) {
	// Verify user exists
	_, err := uc.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Process text into fragments
	totalLines, fragments := uc.textProcessor.ProcessText(input.Content)
	if totalLines == 0 {
		return nil, domain.ErrInvalidTextInfo
	}

	fragmentSize := uc.textProcessor.FragmentSize
	fragmentCount := len(fragments)

	// Generate IDs
	textID := domain.TextID(fmt.Sprintf("text_%d", time.Now().UnixNano()))
	now := time.Now()

	// Create TextInfo
	textInfo, err := domain.NewTextInfo(
		textID,
		input.UserID,
		input.Title,
		totalLines,
		fragmentSize,
		fragmentCount,
		now,
	)
	if err != nil {
		return nil, err
	}

	// Store TextInfo
	if err := uc.textRepo.CreateTextInfo(ctx, textInfo); err != nil {
		return nil, fmt.Errorf("failed to create text info: %w", err)
	}

	// Create and store fragments
	for idx, fragmentLines := range fragments {
		fragmentID := domain.TextFragmentID(fmt.Sprintf("%s_frag_%d", textID, idx))
		fragment, err := domain.NewTextFragment(fragmentID, textID, idx, fragmentLines)
		if err != nil {
			return nil, fmt.Errorf("failed to create fragment %d: %w", idx, err)
		}

		if err := uc.textRepo.CreateFragment(ctx, fragment); err != nil {
			return nil, fmt.Errorf("failed to store fragment %d: %w", idx, err)
		}
	}

	return &CreateTextOutput{TextInfo: textInfo}, nil
}
