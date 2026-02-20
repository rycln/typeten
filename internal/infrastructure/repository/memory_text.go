package repository

import (
	"context"
	"fmt"
	"sync"
	"typeten/internal/domain"
	"typeten/internal/repository"
)

// MemoryTextRepository is an in-memory implementation of TextRepository.
type MemoryTextRepository struct {
	mu        sync.RWMutex
	texts     map[domain.TextID]*domain.TextInfo
	byUser    map[domain.UserID][]*domain.TextInfo
	fragments map[domain.TextFragmentID]*domain.TextFragment
	byTextID  map[domain.TextID][]*domain.TextFragment
}

// NewMemoryTextRepository creates a new in-memory text repository.
func NewMemoryTextRepository() repository.TextRepository {
	return &MemoryTextRepository{
		texts:     make(map[domain.TextID]*domain.TextInfo),
		byUser:    make(map[domain.UserID][]*domain.TextInfo),
		fragments: make(map[domain.TextFragmentID]*domain.TextFragment),
		byTextID:  make(map[domain.TextID][]*domain.TextFragment),
	}
}

func (r *MemoryTextRepository) CreateTextInfo(ctx context.Context, info *domain.TextInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.texts[info.ID]; exists {
		return fmt.Errorf("text already exists")
	}
	
	r.texts[info.ID] = info
	r.byUser[info.UserID] = append(r.byUser[info.UserID], info)
	return nil
}

func (r *MemoryTextRepository) GetTextInfo(ctx context.Context, id domain.TextID) (*domain.TextInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	info, exists := r.texts[id]
	if !exists {
		return nil, fmt.Errorf("text not found")
	}
	return info, nil
}

func (r *MemoryTextRepository) ListByUserID(ctx context.Context, userID domain.UserID) ([]*domain.TextInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	texts := r.byUser[userID]
	if texts == nil {
		return []*domain.TextInfo{}, nil
	}
	
	// Return a copy
	result := make([]*domain.TextInfo, len(texts))
	copy(result, texts)
	return result, nil
}

func (r *MemoryTextRepository) CreateFragment(ctx context.Context, fragment *domain.TextFragment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.fragments[fragment.ID]; exists {
		return fmt.Errorf("fragment already exists")
	}
	
	r.fragments[fragment.ID] = fragment
	r.byTextID[fragment.TextID] = append(r.byTextID[fragment.TextID], fragment)
	return nil
}

func (r *MemoryTextRepository) GetFragment(ctx context.Context, id domain.TextFragmentID) (*domain.TextFragment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	fragment, exists := r.fragments[id]
	if !exists {
		return nil, fmt.Errorf("fragment not found")
	}
	return fragment, nil
}

func (r *MemoryTextRepository) GetFragmentsByTextID(ctx context.Context, textID domain.TextID) ([]*domain.TextFragment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	fragments := r.byTextID[textID]
	if fragments == nil {
		return []*domain.TextFragment{}, nil
	}
	
	// Return a copy
	result := make([]*domain.TextFragment, len(fragments))
	copy(result, fragments)
	return result, nil
}
