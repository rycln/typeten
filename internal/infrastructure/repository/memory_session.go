package repository

import (
	"context"
	"fmt"
	"sync"
	"typeten/internal/domain"
	"typeten/internal/repository"
)

// MemorySessionRepository is an in-memory implementation of SessionRepository.
type MemorySessionRepository struct {
	mu       sync.RWMutex
	sessions map[domain.SessionID]*domain.Session
	byUser   map[domain.UserID][]*domain.Session
}

// NewMemorySessionRepository creates a new in-memory session repository.
func NewMemorySessionRepository() repository.SessionRepository {
	return &MemorySessionRepository{
		sessions: make(map[domain.SessionID]*domain.Session),
		byUser:   make(map[domain.UserID][]*domain.Session),
	}
}

func (r *MemorySessionRepository) Create(ctx context.Context, session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.sessions[session.ID]; exists {
		return fmt.Errorf("session already exists")
	}
	
	r.sessions[session.ID] = session
	r.byUser[session.UserID] = append(r.byUser[session.UserID], session)
	return nil
}

func (r *MemorySessionRepository) GetByID(ctx context.Context, id domain.SessionID) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	session, exists := r.sessions[id]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}

func (r *MemorySessionRepository) Update(ctx context.Context, session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.sessions[session.ID]; !exists {
		return fmt.Errorf("session not found")
	}
	
	r.sessions[session.ID] = session
	return nil
}

func (r *MemorySessionRepository) ListByUserID(ctx context.Context, userID domain.UserID) ([]*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	sessions := r.byUser[userID]
	if sessions == nil {
		return []*domain.Session{}, nil
	}
	
	// Return a copy
	result := make([]*domain.Session, len(sessions))
	copy(result, sessions)
	return result, nil
}
