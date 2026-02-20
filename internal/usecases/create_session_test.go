package usecases

import (
	"context"
	"testing"
	"time"
	"typeten/internal/domain"
)

func TestCreateSessionUseCase_Execute(t *testing.T) {
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
	textInfo, err := domain.NewTextInfo("text_1", user.ID, "Test Text", 10, 5, 2, now)
	if err != nil {
		t.Fatalf("Failed to create text info: %v", err)
	}
	if err := textRepo.CreateTextInfo(ctx, textInfo); err != nil {
		t.Fatalf("Failed to store text info: %v", err)
	}

	sessionRepo := NewMockSessionRepository()
	useCase := NewCreateSessionUseCase(sessionRepo, textRepo, userRepo)

	tests := []struct {
		name    string
		input   CreateSessionInput
		wantErr bool
	}{
		{
			name: "valid session",
			input: CreateSessionInput{
				UserID: user.ID,
				TextID: textInfo.ID,
			},
			wantErr: false,
		},
		{
			name: "non-existent user",
			input: CreateSessionInput{
				UserID: "nonexistent",
				TextID: textInfo.ID,
			},
			wantErr: true,
		},
		{
			name: "non-existent text",
			input: CreateSessionInput{
				UserID: user.ID,
				TextID: "nonexistent",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset session repo for each test
			sessionRepo = NewMockSessionRepository()
			useCase = NewCreateSessionUseCase(sessionRepo, textRepo, userRepo)

			output, err := useCase.Execute(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if output == nil {
					t.Fatal("Execute() returned nil output")
				}
				if output.Session == nil {
					t.Fatal("Execute() returned nil Session")
				}
				if output.Session.UserID != tt.input.UserID {
					t.Errorf("Execute() UserID = %v, want %v", output.Session.UserID, tt.input.UserID)
				}
				if output.Session.TextID != tt.input.TextID {
					t.Errorf("Execute() TextID = %v, want %v", output.Session.TextID, tt.input.TextID)
				}
				// Verify session was stored
				stored, err := sessionRepo.GetByID(ctx, output.Session.ID)
				if err != nil {
					t.Errorf("Failed to get stored session: %v", err)
				}
				if stored == nil {
					t.Error("Session was not stored")
				}
			}
		})
	}
}
