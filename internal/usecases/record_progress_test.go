package usecases

import (
	"context"
	"testing"
	"time"
	"typeten/internal/domain"
)

func TestRecordProgressUseCase_Execute(t *testing.T) {
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

	useCase := NewRecordProgressUseCase(sessionRepo)

	tests := []struct {
		name    string
		input   RecordProgressInput
		wantErr bool
	}{
		{
			name: "valid progress",
			input: RecordProgressInput{
				SessionID:       "session_1",
				AccuracyPercent: 95.5,
				WPM:             45.2,
			},
			wantErr: false,
		},
		{
			name: "non-existent session",
			input: RecordProgressInput{
				SessionID:       "nonexistent",
				AccuracyPercent: 95.0,
				WPM:             45.0,
			},
			wantErr: true,
		},
		{
			name: "invalid accuracy",
			input: RecordProgressInput{
				SessionID:       "session_1",
				AccuracyPercent: 150.0,
				WPM:             45.0,
			},
			wantErr: true,
		},
		{
			name: "negative WPM",
			input: RecordProgressInput{
				SessionID:       "session_1",
				AccuracyPercent: 95.0,
				WPM:             -10.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset session for each test
			session, err := domain.NewSession("session_1", "user_1", "text_1", now)
			if err != nil {
				t.Fatalf("Failed to create session: %v", err)
			}
			sessionRepo = NewMockSessionRepository()
			if err := sessionRepo.Create(ctx, session); err != nil {
				t.Fatalf("Failed to store session: %v", err)
			}
			useCase = NewRecordProgressUseCase(sessionRepo)

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
				if output.Session.CompletedLines != 1 {
					t.Errorf("Execute() CompletedLines = %v, want 1", output.Session.CompletedLines)
				}
			}
		})
	}
}
