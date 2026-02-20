package usecases

import (
	"context"
	"testing"
	"time"
	"typeten/internal/domain"
)

func TestGetSessionUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	sessionRepo := NewMockSessionRepository()
	session, err := domain.NewSession("session_1", "user_1", "text_1", now)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
	if err := sessionRepo.Create(ctx, session); err != nil {
		t.Fatalf("Failed to store session: %v", err)
	}

	useCase := NewGetSessionUseCase(sessionRepo)

	tests := []struct {
		name    string
		input   GetSessionInput
		wantErr bool
	}{
		{
			name: "valid session",
			input: GetSessionInput{
				SessionID: "session_1",
			},
			wantErr: false,
		},
		{
			name: "non-existent session",
			input: GetSessionInput{
				SessionID: "nonexistent",
			},
			wantErr: true,
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
				if output.Session == nil {
					t.Fatal("Execute() returned nil Session")
				}
				if output.Session.ID != tt.input.SessionID {
					t.Errorf("Execute() SessionID = %v, want %v", output.Session.ID, tt.input.SessionID)
				}
			}
		})
	}
}
