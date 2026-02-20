package repository

import (
	"context"
	"testing"
	"time"
	"typeten/internal/domain"
)

func TestMemorySessionRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewMemorySessionRepository()

	now := time.Now()
	session1, err := domain.NewSession("session_1", "user_1", "text_1", now)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	session2, err := domain.NewSession("session_2", "user_1", "text_2", now)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	t.Run("Create and GetByID", func(t *testing.T) {
		if err := repo.Create(ctx, session1); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.GetByID(ctx, session1.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.ID != session1.ID {
			t.Errorf("GetByID() ID = %v, want %v", got.ID, session1.ID)
		}
		if got.UserID != session1.UserID {
			t.Errorf("GetByID() UserID = %v, want %v", got.UserID, session1.UserID)
		}
	})

	t.Run("GetByID non-existent", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "nonexistent")
		if err == nil {
			t.Error("GetByID() expected error for non-existent session")
		}
	})

	t.Run("Create duplicate", func(t *testing.T) {
		err := repo.Create(ctx, session1)
		if err == nil {
			t.Error("Create() expected error for duplicate session")
		}
	})

	t.Run("Update", func(t *testing.T) {
		if err := repo.Create(ctx, session2); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		session2.CompletedLines = 5
		if err := repo.Update(ctx, session2); err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		got, err := repo.GetByID(ctx, session2.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.CompletedLines != 5 {
			t.Errorf("Update() CompletedLines = %v, want 5", got.CompletedLines)
		}
	})

	t.Run("Update non-existent", func(t *testing.T) {
		nonExistent, _ := domain.NewSession("nonexistent", "user_1", "text_1", now)
		err := repo.Update(ctx, nonExistent)
		if err == nil {
			t.Error("Update() expected error for non-existent session")
		}
	})

	t.Run("ListByUserID", func(t *testing.T) {
		sessions, err := repo.ListByUserID(ctx, "user_1")
		if err != nil {
			t.Fatalf("ListByUserID() error = %v", err)
		}
		if len(sessions) < 2 {
			t.Errorf("ListByUserID() length = %v, want at least 2", len(sessions))
		}
	})

	t.Run("ListByUserID empty", func(t *testing.T) {
		sessions, err := repo.ListByUserID(ctx, "nonexistent")
		if err != nil {
			t.Fatalf("ListByUserID() error = %v", err)
		}
		if len(sessions) != 0 {
			t.Errorf("ListByUserID() length = %v, want 0", len(sessions))
		}
	})
}
