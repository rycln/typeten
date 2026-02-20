package usecases

import (
	"context"
	"testing"
	"time"
	"typeten/internal/domain"
)

func TestListTextsUseCase_Execute(t *testing.T) {
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
	text1, err := domain.NewTextInfo("text_1", user.ID, "Text 1", 10, 5, 2, now)
	if err != nil {
		t.Fatalf("Failed to create text: %v", err)
	}
	if err := textRepo.CreateTextInfo(ctx, text1); err != nil {
		t.Fatalf("Failed to store text: %v", err)
	}

	text2, err := domain.NewTextInfo("text_2", user.ID, "Text 2", 20, 10, 2, now)
	if err != nil {
		t.Fatalf("Failed to create text: %v", err)
	}
	if err := textRepo.CreateTextInfo(ctx, text2); err != nil {
		t.Fatalf("Failed to store text: %v", err)
	}

	useCase := NewListTextsUseCase(textRepo, userRepo)

	tests := []struct {
		name    string
		input   ListTextsInput
		wantErr bool
		wantLen int
	}{
		{
			name: "list user texts",
			input: ListTextsInput{
				UserID: user.ID,
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "non-existent user",
			input: ListTextsInput{
				UserID: "nonexistent",
			},
			wantErr: true,
			wantLen: 0,
		},
		{
			name: "user with no texts",
			input: ListTextsInput{
				UserID: "user_no_texts",
			},
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := useCase.Execute(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if output == nil {
					t.Fatal("Execute() returned nil output")
				}
				if len(output.Texts) != tt.wantLen {
					t.Errorf("Execute() Texts length = %v, want %v", len(output.Texts), tt.wantLen)
				}
			}
		})
	}
}
