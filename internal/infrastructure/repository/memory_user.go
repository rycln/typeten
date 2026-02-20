package repository

import (
	"context"
	"fmt"
	"sync"
	"typeten/internal/domain"
	"typeten/internal/repository"
)

// MemoryUserRepository is an in-memory implementation of UserRepository.
type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[domain.UserID]*domain.User
	byEmail map[string]*domain.User
}

// NewMemoryUserRepository creates a new in-memory user repository.
func NewMemoryUserRepository() repository.UserRepository {
	return &MemoryUserRepository{
		users:   make(map[domain.UserID]*domain.User),
		byEmail: make(map[string]*domain.User),
	}
}

func (r *MemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("user already exists")
	}
	
	r.users[user.ID] = user
	r.byEmail[user.Email] = user
	return nil
}

func (r *MemoryUserRepository) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *MemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	user, exists := r.byEmail[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}
