package usecases

import (
	"context"
	"fmt"
	"sort"
	"typeten/internal/domain"
	"typeten/internal/repository"
)

// GetTextFragmentsUseCase handles retrieving fragments for a text.
type GetTextFragmentsUseCase struct {
	textRepo repository.TextRepository
}

// NewGetTextFragmentsUseCase creates a new GetTextFragmentsUseCase.
func NewGetTextFragmentsUseCase(textRepo repository.TextRepository) *GetTextFragmentsUseCase {
	return &GetTextFragmentsUseCase{
		textRepo: textRepo,
	}
}

// GetTextFragmentsInput represents the input for getting fragments.
type GetTextFragmentsInput struct {
	TextID domain.TextID
}

// GetTextFragmentsOutput represents the result of getting fragments.
type GetTextFragmentsOutput struct {
	Fragments []*domain.TextFragment
}

// Execute retrieves all fragments for a text, ordered by FragmentIdx.
func (uc *GetTextFragmentsUseCase) Execute(ctx context.Context, input GetTextFragmentsInput) (*GetTextFragmentsOutput, error) {
	// Verify text exists
	_, err := uc.textRepo.GetTextInfo(ctx, input.TextID)
	if err != nil {
		return nil, fmt.Errorf("text not found: %w", err)
	}
	
	fragments, err := uc.textRepo.GetFragmentsByTextID(ctx, input.TextID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fragments: %w", err)
	}
	
	// Sort by FragmentIdx
	sort.Slice(fragments, func(i, j int) bool {
		return fragments[i].FragmentIdx < fragments[j].FragmentIdx
	})
	
	return &GetTextFragmentsOutput{Fragments: fragments}, nil
}
