package domain

import (
	"strings"
	"time"
)

// TextID identifies a text. TextFragmentID identifies a single fragment of that text.
type (
	TextID         string
	TextFragmentID string
)

// TextInfo holds metadata for a single text: ownership, title, and layout
// (total lines, fragment size and count). The actual content is stored as
// one or more TextFragment values referenced by TextID.
type TextInfo struct {
	ID            TextID
	UserID        UserID
	Title         string
	TotalLines    int
	FragmentSize  int
	FragmentCount int
	CreatedAt     time.Time
}

// NewTextInfo creates a TextInfo after validating IDs, title, and line/fragment counts.
// Returns ErrInvalidTextInfo if any field is invalid.
func NewTextInfo(id TextID, userID UserID, title string, totalLines, fragmentSize, fragmentCount int, createdAt time.Time) (*TextInfo, error) {
	if err := validateTextID(id); err != nil {
		return nil, err
	}
	if err := validateUserID(userID); err != nil {
		return nil, err
	}
	if strings.TrimSpace(title) == "" {
		return nil, ErrInvalidTextInfo
	}
	if totalLines <= 0 || fragmentSize <= 0 || fragmentCount <= 0 {
		return nil, ErrInvalidTextInfo
	}
	return &TextInfo{
		ID:            id,
		UserID:        userID,
		Title:         strings.TrimSpace(title),
		TotalLines:    totalLines,
		FragmentSize:  fragmentSize,
		FragmentCount: fragmentCount,
		CreatedAt:     createdAt,
	}, nil
}

// TextFragment is one chunk of a text's content (Lines), with FragmentIdx
// indicating its order. The full text is the ordered set of fragments for a given TextID.
type TextFragment struct {
	ID          TextFragmentID
	TextID      TextID
	FragmentIdx int
	lines       []string // copied on construction; use Lines() to read a copy
}

// NewTextFragment creates a TextFragment, copying lines so the slice cannot be mutated by callers.
// Returns ErrInvalidFragment if id, textID are empty, fragmentIdx < 0, or lines is empty.
func NewTextFragment(id TextFragmentID, textID TextID, fragmentIdx int, lines []string) (*TextFragment, error) {
	if err := validateTextFragmentID(id); err != nil {
		return nil, err
	}
	if err := validateTextID(textID); err != nil {
		return nil, err
	}
	if fragmentIdx < 0 {
		return nil, ErrInvalidFragment
	}
	if len(lines) == 0 {
		return nil, ErrInvalidFragment
	}
	cp := make([]string, len(lines))
	copy(cp, lines)
	return &TextFragment{
		ID:          id,
		TextID:      textID,
		FragmentIdx: fragmentIdx,
		lines:       cp,
	}, nil
}

// Lines returns a copy of the fragment's lines. Callers must not rely on mutating the returned slice.
func (f *TextFragment) Lines() []string {
	if f == nil || len(f.lines) == 0 {
		return nil
	}
	out := make([]string, len(f.lines))
	copy(out, f.lines)
	return out
}

func validateTextID(id TextID) error {
	if strings.TrimSpace(string(id)) == "" {
		return ErrInvalidID
	}
	return nil
}

func validateTextFragmentID(id TextFragmentID) error {
	if strings.TrimSpace(string(id)) == "" {
		return ErrInvalidID
	}
	return nil
}
