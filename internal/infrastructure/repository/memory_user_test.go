package repository

import (
	"context"
	"testing"
	"time"
	"typeten/internal/domain"
)

func TestMemoryUserRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryUserRepository()

	now := time.Now()
	user1, err := domain.NewUser("user_1", "test1@example.com", "user1", now)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user2, err := domain.NewUser("user_2", "test2@example.com", "user2", now)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	t.Run("Create and GetByID", func(t *testing.T) {
		if err := repo.Create(ctx, user1); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.GetByID(ctx, user1.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.ID != user1.ID {
			t.Errorf("GetByID() ID = %v, want %v", got.ID, user1.ID)
		}
		if got.Email != user1.Email {
			t.Errorf("GetByID() Email = %v, want %v", got.Email, user1.Email)
		}
	})

	t.Run("GetByID non-existent", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "nonexistent")
		if err == nil {
			t.Error("GetByID() expected error for non-existent user")
		}
	})

	t.Run("Create duplicate", func(t *testing.T) {
		err := repo.Create(ctx, user1)
		if err == nil {
			t.Error("Create() expected error for duplicate user")
		}
	})

	t.Run("GetByEmail", func(t *testing.T) {
		if err := repo.Create(ctx, user2); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.GetByEmail(ctx, user2.Email)
		if err != nil {
			t.Fatalf("GetByEmail() error = %v", err)
		}
		if got.ID != user2.ID {
			t.Errorf("GetByEmail() ID = %v, want %v", got.ID, user2.ID)
		}
	})

	t.Run("GetByEmail non-existent", func(t *testing.T) {
		_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
		if err == nil {
			t.Error("GetByEmail() expected error for non-existent email")
		}
	})
}
