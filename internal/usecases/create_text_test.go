package usecases

import (
	"context"
	"testing"
	"time"
	"typeten/internal/domain"
)

func TestCreateTextUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userRepo := NewMockUserRepository()
	user, err := domain.NewUser("user_1", "test@example.com", "testuser", now)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to store user: %v", err)
	}

	textRepo := NewMockTextRepository()
	useCase := NewCreateTextUseCase(textRepo, userRepo, 5)

	tests := []struct {
		name    string
		input   CreateTextInput
		wantErr bool
	}{
		{
			name: "valid text",
			input: CreateTextInput{
				UserID:  user.ID,
				Title:   "Test Text",
				Content: "line1\nline2\nline3",
			},
			wantErr: false,
		},
		{
			name: "whitespace only content",
			input: CreateTextInput{
				UserID:  user.ID,
				Title:   "Empty Text",
				Content: "   \n  \n  ",
			},
			wantErr: true,
		},
		{
			name: "non-existent user",
			input: CreateTextInput{
				UserID:  "nonexistent",
				Title:   "Test",
				Content: "line1",
			},
			wantErr: true,
		},
		{
			name: "multi-line text",
			input: CreateTextInput{
				UserID:  user.ID,
				Title:   "Multi-line",
				Content: "line1\nline2\nline3\nline4\nline5\nline6",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset text repo for each test
			textRepo = NewMockTextRepository()
			useCase = NewCreateTextUseCase(textRepo, userRepo, 5)

			output, err := useCase.Execute(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if output == nil {
					t.Fatal("Execute() returned nil output")
				}
				if output.TextInfo == nil {
					t.Fatal("Execute() returned nil TextInfo")
				}
				if output.TextInfo.Title != tt.input.Title {
					t.Errorf("Execute() Title = %v, want %v", output.TextInfo.Title, tt.input.Title)
				}
				// Verify fragments were created
				frags, err := textRepo.GetFragmentsByTextID(ctx, output.TextInfo.ID)
				if err != nil {
					t.Errorf("Failed to get fragments: %v", err)
				}
				if len(frags) == 0 && tt.input.Content != "" {
					t.Error("Execute() created no fragments")
				}
			}
		})
	}
}
