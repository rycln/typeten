package usecases

import (
	"context"
	"testing"
	"time"
	"typeten/internal/domain"
)

func TestGetTextFragmentsUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	textRepo := NewMockTextRepository()
	textInfo, err := domain.NewTextInfo("text_1", "user_1", "Test Text", 10, 5, 2, now)
	if err != nil {
		t.Fatalf("Failed to create text info: %v", err)
	}
	if err := textRepo.CreateTextInfo(ctx, textInfo); err != nil {
		t.Fatalf("Failed to store text info: %v", err)
	}

	frag1, err := domain.NewTextFragment("frag_1", textInfo.ID, 0, []string{"line1", "line2"})
	if err != nil {
		t.Fatalf("Failed to create fragment: %v", err)
	}
	if err := textRepo.CreateFragment(ctx, frag1); err != nil {
		t.Fatalf("Failed to store fragment: %v", err)
	}

	frag2, err := domain.NewTextFragment("frag_2", textInfo.ID, 1, []string{"line3", "line4"})
	if err != nil {
		t.Fatalf("Failed to create fragment: %v", err)
	}
	if err := textRepo.CreateFragment(ctx, frag2); err != nil {
		t.Fatalf("Failed to store fragment: %v", err)
	}

	useCase := NewGetTextFragmentsUseCase(textRepo)

	tests := []struct {
		name    string
		input   GetTextFragmentsInput
		wantErr bool
		wantLen int
	}{
		{
			name: "valid text",
			input: GetTextFragmentsInput{
				TextID: textInfo.ID,
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "non-existent text",
			input: GetTextFragmentsInput{
				TextID: "nonexistent",
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
				if len(output.Fragments) != tt.wantLen {
					t.Errorf("Execute() Fragments length = %v, want %v", len(output.Fragments), tt.wantLen)
				}
				// Verify fragments are sorted by FragmentIdx
				for i := 1; i < len(output.Fragments); i++ {
					if output.Fragments[i].FragmentIdx <= output.Fragments[i-1].FragmentIdx {
						t.Error("Fragments are not sorted by FragmentIdx")
					}
				}
			}
		})
	}
}
