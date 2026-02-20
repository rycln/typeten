package usecases

import (
	"context"
	"fmt"
	"typeten/internal/domain"
)

// MockUserRepository is a mock implementation of UserRepository for testing.
type MockUserRepository struct {
	users   map[domain.UserID]*domain.User
	byEmail map[string]*domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:   make(map[domain.UserID]*domain.User),
		byEmail: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if _, exists := m.users[user.ID]; exists {
		return fmt.Errorf("user already exists")
	}
	m.users[user.ID] = user
	m.byEmail[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, exists := m.byEmail[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// MockTextRepository is a mock implementation of TextRepository for testing.
type MockTextRepository struct {
	texts     map[domain.TextID]*domain.TextInfo
	byUser    map[domain.UserID][]*domain.TextInfo
	fragments map[domain.TextFragmentID]*domain.TextFragment
	byTextID  map[domain.TextID][]*domain.TextFragment
}

func NewMockTextRepository() *MockTextRepository {
	return &MockTextRepository{
		texts:     make(map[domain.TextID]*domain.TextInfo),
		byUser:    make(map[domain.UserID][]*domain.TextInfo),
		fragments: make(map[domain.TextFragmentID]*domain.TextFragment),
		byTextID:  make(map[domain.TextID][]*domain.TextFragment),
	}
}

func (m *MockTextRepository) CreateTextInfo(ctx context.Context, info *domain.TextInfo) error {
	if _, exists := m.texts[info.ID]; exists {
		return fmt.Errorf("text already exists")
	}
	m.texts[info.ID] = info
	m.byUser[info.UserID] = append(m.byUser[info.UserID], info)
	return nil
}

func (m *MockTextRepository) GetTextInfo(ctx context.Context, id domain.TextID) (*domain.TextInfo, error) {
	info, exists := m.texts[id]
	if !exists {
		return nil, fmt.Errorf("text not found")
	}
	return info, nil
}

func (m *MockTextRepository) ListByUserID(ctx context.Context, userID domain.UserID) ([]*domain.TextInfo, error) {
	texts := m.byUser[userID]
	if texts == nil {
		return []*domain.TextInfo{}, nil
	}
	result := make([]*domain.TextInfo, len(texts))
	copy(result, texts)
	return result, nil
}

func (m *MockTextRepository) CreateFragment(ctx context.Context, fragment *domain.TextFragment) error {
	if _, exists := m.fragments[fragment.ID]; exists {
		return fmt.Errorf("fragment already exists")
	}
	m.fragments[fragment.ID] = fragment
	m.byTextID[fragment.TextID] = append(m.byTextID[fragment.TextID], fragment)
	return nil
}

func (m *MockTextRepository) GetFragment(ctx context.Context, id domain.TextFragmentID) (*domain.TextFragment, error) {
	fragment, exists := m.fragments[id]
	if !exists {
		return nil, fmt.Errorf("fragment not found")
	}
	return fragment, nil
}

func (m *MockTextRepository) GetFragmentsByTextID(ctx context.Context, textID domain.TextID) ([]*domain.TextFragment, error) {
	fragments := m.byTextID[textID]
	if fragments == nil {
		return []*domain.TextFragment{}, nil
	}
	result := make([]*domain.TextFragment, len(fragments))
	copy(result, fragments)
	return result, nil
}

// MockSessionRepository is a mock implementation of SessionRepository for testing.
type MockSessionRepository struct {
	sessions map[domain.SessionID]*domain.Session
	byUser   map[domain.UserID][]*domain.Session
}

func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{
		sessions: make(map[domain.SessionID]*domain.Session),
		byUser:   make(map[domain.UserID][]*domain.Session),
	}
}

func (m *MockSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	if _, exists := m.sessions[session.ID]; exists {
		return fmt.Errorf("session already exists")
	}
	m.sessions[session.ID] = session
	m.byUser[session.UserID] = append(m.byUser[session.UserID], session)
	return nil
}

func (m *MockSessionRepository) GetByID(ctx context.Context, id domain.SessionID) (*domain.Session, error) {
	session, exists := m.sessions[id]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}

func (m *MockSessionRepository) Update(ctx context.Context, session *domain.Session) error {
	if _, exists := m.sessions[session.ID]; !exists {
		return fmt.Errorf("session not found")
	}
	m.sessions[session.ID] = session
	return nil
}

func (m *MockSessionRepository) ListByUserID(ctx context.Context, userID domain.UserID) ([]*domain.Session, error) {
	sessions := m.byUser[userID]
	if sessions == nil {
		return []*domain.Session{}, nil
	}
	result := make([]*domain.Session, len(sessions))
	copy(result, sessions)
	return result, nil
}
