package repository

import (
	"context"
	"testing"
	"time"
	"typeten/internal/domain"
)

func TestMemoryTextRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryTextRepository()

	now := time.Now()
	userID := domain.UserID("user_1")
	textID := domain.TextID("text_1")

	textInfo, err := domain.NewTextInfo(textID, userID, "Test Text", 10, 5, 2, now)
	if err != nil {
		t.Fatalf("Failed to create text info: %v", err)
	}

	frag1, err := domain.NewTextFragment("frag_1", textID, 0, []string{"line1", "line2"})
	if err != nil {
		t.Fatalf("Failed to create fragment: %v", err)
	}

	frag2, err := domain.NewTextFragment("frag_2", textID, 1, []string{"line3", "line4"})
	if err != nil {
		t.Fatalf("Failed to create fragment: %v", err)
	}

	t.Run("CreateTextInfo and GetTextInfo", func(t *testing.T) {
		if err := repo.CreateTextInfo(ctx, textInfo); err != nil {
			t.Fatalf("CreateTextInfo() error = %v", err)
		}

		got, err := repo.GetTextInfo(ctx, textID)
		if err != nil {
			t.Fatalf("GetTextInfo() error = %v", err)
		}
		if got.ID != textID {
			t.Errorf("GetTextInfo() ID = %v, want %v", got.ID, textID)
		}
		if got.Title != textInfo.Title {
			t.Errorf("GetTextInfo() Title = %v, want %v", got.Title, textInfo.Title)
		}
	})

	t.Run("GetTextInfo non-existent", func(t *testing.T) {
		_, err := repo.GetTextInfo(ctx, "nonexistent")
		if err == nil {
			t.Error("GetTextInfo() expected error for non-existent text")
		}
	})

	t.Run("CreateTextInfo duplicate", func(t *testing.T) {
		err := repo.CreateTextInfo(ctx, textInfo)
		if err == nil {
			t.Error("CreateTextInfo() expected error for duplicate text")
		}
	})

	t.Run("ListByUserID", func(t *testing.T) {
		texts, err := repo.ListByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("ListByUserID() error = %v", err)
		}
		if len(texts) == 0 {
			t.Error("ListByUserID() returned empty list")
		}
		if texts[0].ID != textID {
			t.Errorf("ListByUserID() ID = %v, want %v", texts[0].ID, textID)
		}
	})

	t.Run("ListByUserID empty", func(t *testing.T) {
		texts, err := repo.ListByUserID(ctx, "nonexistent")
		if err != nil {
			t.Fatalf("ListByUserID() error = %v", err)
		}
		if len(texts) != 0 {
			t.Errorf("ListByUserID() length = %v, want 0", len(texts))
		}
	})

	t.Run("CreateFragment and GetFragment", func(t *testing.T) {
		if err := repo.CreateFragment(ctx, frag1); err != nil {
			t.Fatalf("CreateFragment() error = %v", err)
		}

		got, err := repo.GetFragment(ctx, frag1.ID)
		if err != nil {
			t.Fatalf("GetFragment() error = %v", err)
		}
		if got.ID != frag1.ID {
			t.Errorf("GetFragment() ID = %v, want %v", got.ID, frag1.ID)
		}
	})

	t.Run("GetFragment non-existent", func(t *testing.T) {
		_, err := repo.GetFragment(ctx, "nonexistent")
		if err == nil {
			t.Error("GetFragment() expected error for non-existent fragment")
		}
	})

	t.Run("GetFragmentsByTextID", func(t *testing.T) {
		if err := repo.CreateFragment(ctx, frag2); err != nil {
			t.Fatalf("CreateFragment() error = %v", err)
		}

		frags, err := repo.GetFragmentsByTextID(ctx, textID)
		if err != nil {
			t.Fatalf("GetFragmentsByTextID() error = %v", err)
		}
		if len(frags) < 2 {
			t.Errorf("GetFragmentsByTextID() length = %v, want at least 2", len(frags))
		}
	})

	t.Run("GetFragmentsByTextID empty", func(t *testing.T) {
		frags, err := repo.GetFragmentsByTextID(ctx, "nonexistent")
		if err != nil {
			t.Fatalf("GetFragmentsByTextID() error = %v", err)
		}
		if len(frags) != 0 {
			t.Errorf("GetFragmentsByTextID() length = %v, want 0", len(frags))
		}
	})
}
