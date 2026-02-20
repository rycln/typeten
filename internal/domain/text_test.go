package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewTextInfo(t *testing.T) {
	now := time.Now()
	validID := TextID("text_123")
	validUserID := UserID("user_123")
	validTitle := "Test Text"
	validTotalLines := 10
	validFragmentSize := 5
	validFragmentCount := 2

	tests := []struct {
		name          string
		id            TextID
		userID        UserID
		title         string
		totalLines    int
		fragmentSize  int
		fragmentCount int
		createdAt     time.Time
		wantErr       error
	}{
		{
			name:          "valid text info",
			id:            validID,
			userID:        validUserID,
			title:         validTitle,
			totalLines:    validTotalLines,
			fragmentSize:  validFragmentSize,
			fragmentCount: validFragmentCount,
			createdAt:     now,
			wantErr:       nil,
		},
		{
			name:          "empty text ID",
			id:            "",
			userID:        validUserID,
			title:         validTitle,
			totalLines:    validTotalLines,
			fragmentSize:  validFragmentSize,
			fragmentCount: validFragmentCount,
			createdAt:     now,
			wantErr:       ErrInvalidID,
		},
		{
			name:          "empty user ID",
			id:            validID,
			userID:        "",
			title:         validTitle,
			totalLines:    validTotalLines,
			fragmentSize:  validFragmentSize,
			fragmentCount: validFragmentCount,
			createdAt:     now,
			wantErr:       ErrInvalidID,
		},
		{
			name:          "empty title",
			id:            validID,
			userID:        validUserID,
			title:         "",
			totalLines:    validTotalLines,
			fragmentSize:  validFragmentSize,
			fragmentCount: validFragmentCount,
			createdAt:     now,
			wantErr:       ErrInvalidTextInfo,
		},
		{
			name:          "whitespace title",
			id:            validID,
			userID:        validUserID,
			title:         "   ",
			totalLines:    validTotalLines,
			fragmentSize:  validFragmentSize,
			fragmentCount: validFragmentCount,
			createdAt:     now,
			wantErr:       ErrInvalidTextInfo,
		},
		{
			name:          "zero total lines",
			id:            validID,
			userID:        validUserID,
			title:         validTitle,
			totalLines:    0,
			fragmentSize:  validFragmentSize,
			fragmentCount: validFragmentCount,
			createdAt:     now,
			wantErr:       ErrInvalidTextInfo,
		},
		{
			name:          "negative total lines",
			id:            validID,
			userID:        validUserID,
			title:         validTitle,
			totalLines:    -1,
			fragmentSize:  validFragmentSize,
			fragmentCount: validFragmentCount,
			createdAt:     now,
			wantErr:       ErrInvalidTextInfo,
		},
		{
			name:          "zero fragment size",
			id:            validID,
			userID:        validUserID,
			title:         validTitle,
			totalLines:    validTotalLines,
			fragmentSize:  0,
			fragmentCount: validFragmentCount,
			createdAt:     now,
			wantErr:       ErrInvalidTextInfo,
		},
		{
			name:          "zero fragment count",
			id:            validID,
			userID:        validUserID,
			title:         validTitle,
			totalLines:    validTotalLines,
			fragmentSize:  validFragmentSize,
			fragmentCount: 0,
			createdAt:     now,
			wantErr:       ErrInvalidTextInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := NewTextInfo(tt.id, tt.userID, tt.title, tt.totalLines, tt.fragmentSize, tt.fragmentCount, tt.createdAt)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewTextInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if info == nil {
					t.Fatal("NewTextInfo() returned nil")
				}
				if info.ID != tt.id {
					t.Errorf("NewTextInfo() ID = %v, want %v", info.ID, tt.id)
				}
				if info.Title != strings.TrimSpace(tt.title) {
					t.Errorf("NewTextInfo() Title = %v, want %v", info.Title, strings.TrimSpace(tt.title))
				}
			}
		})
	}
}

func TestNewTextFragment(t *testing.T) {
	validID := TextFragmentID("frag_123")
	validTextID := TextID("text_123")
	validLines := []string{"line 1", "line 2"}

	tests := []struct {
		name       string
		id         TextFragmentID
		textID     TextID
		fragmentIdx int
		lines      []string
		wantErr    error
	}{
		{
			name:        "valid fragment",
			id:          validID,
			textID:      validTextID,
			fragmentIdx: 0,
			lines:       validLines,
			wantErr:     nil,
		},
		{
			name:        "empty fragment ID",
			id:          "",
			textID:      validTextID,
			fragmentIdx: 0,
			lines:       validLines,
			wantErr:     ErrInvalidID,
		},
		{
			name:        "empty text ID",
			id:          validID,
			textID:      "",
			fragmentIdx: 0,
			lines:       validLines,
			wantErr:     ErrInvalidID,
		},
		{
			name:        "negative fragment index",
			id:          validID,
			textID:      validTextID,
			fragmentIdx: -1,
			lines:       validLines,
			wantErr:     ErrInvalidFragment,
		},
		{
			name:        "empty lines",
			id:          validID,
			textID:      validTextID,
			fragmentIdx: 0,
			lines:       []string{},
			wantErr:     ErrInvalidFragment,
		},
		{
			name:        "nil lines",
			id:          validID,
			textID:      validTextID,
			fragmentIdx: 0,
			lines:       nil,
			wantErr:     ErrInvalidFragment,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frag, err := NewTextFragment(tt.id, tt.textID, tt.fragmentIdx, tt.lines)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewTextFragment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if frag == nil {
					t.Fatal("NewTextFragment() returned nil")
				}
				if len(frag.Lines()) != len(tt.lines) {
					t.Errorf("NewTextFragment() Lines() length = %v, want %v", len(frag.Lines()), len(tt.lines))
				}
				// Verify lines are copied (not same slice)
				if len(tt.lines) > 0 && &frag.Lines()[0] == &tt.lines[0] {
					t.Error("NewTextFragment() Lines() returned same slice, should be a copy")
				}
			}
		})
	}
}

func TestTextFragment_Lines(t *testing.T) {
	frag, err := NewTextFragment("frag_1", "text_1", 0, []string{"line 1", "line 2"})
	if err != nil {
		t.Fatalf("NewTextFragment() error = %v", err)
	}

	lines1 := frag.Lines()
	lines2 := frag.Lines()

	// Verify each call returns a copy
	if &lines1[0] == &lines2[0] {
		t.Error("Lines() returned same slice on multiple calls, should return copies")
	}

	// Verify modifying returned slice doesn't affect fragment
	lines1[0] = "modified"
	if frag.Lines()[0] == "modified" {
		t.Error("Modifying returned Lines() affected fragment internal state")
	}

	// Test nil fragment
	var nilFrag *TextFragment
	if nilFrag.Lines() != nil {
		t.Error("nil fragment Lines() should return nil")
	}
}
