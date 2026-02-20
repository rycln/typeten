package repository

import (
	"context"
	"typeten/internal/domain"
)

// UserRepository defines operations for user persistence.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// TextRepository defines operations for text persistence.
type TextRepository interface {
	CreateTextInfo(ctx context.Context, info *domain.TextInfo) error
	GetTextInfo(ctx context.Context, id domain.TextID) (*domain.TextInfo, error)
	ListByUserID(ctx context.Context, userID domain.UserID) ([]*domain.TextInfo, error)
	
	CreateFragment(ctx context.Context, fragment *domain.TextFragment) error
	GetFragment(ctx context.Context, id domain.TextFragmentID) (*domain.TextFragment, error)
	GetFragmentsByTextID(ctx context.Context, textID domain.TextID) ([]*domain.TextFragment, error)
}

// SessionRepository defines operations for session persistence.
type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, id domain.SessionID) (*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) error
	ListByUserID(ctx context.Context, userID domain.UserID) ([]*domain.Session, error)
}
