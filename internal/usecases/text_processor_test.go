package usecases

import "testing"

func TestTextProcessor_ProcessText(t *testing.T) {
	tests := []struct {
		name         string
		fragmentSize int
		text         string
		wantLines    int
		wantFrags    int
	}{
		{
			name:         "single fragment",
			fragmentSize: 10,
			text:         "line1\nline2\nline3",
			wantLines:    3,
			wantFrags:    1,
		},
		{
			name:         "multiple fragments",
			fragmentSize: 2,
			text:         "line1\nline2\nline3\nline4\nline5",
			wantLines:    5,
			wantFrags:    3, // 2+2+1
		},
		{
			name:         "empty text",
			fragmentSize: 10,
			text:         "",
			wantLines:    0,
			wantFrags:    0,
		},
		{
			name:         "whitespace only text",
			fragmentSize: 10,
			text:         "   \n  \n  ",
			wantLines:    0,
			wantFrags:    0,
		},
		{
			name:         "whitespace only",
			fragmentSize: 10,
			text:         "line1\nline2\nline3",
			wantLines:    3,
			wantFrags:    1,
		},
		{
			name:         "single line",
			fragmentSize: 10,
			text:         "single line",
			wantLines:    1,
			wantFrags:    1,
		},
		{
			name:         "exact fragment size",
			fragmentSize: 3,
			text:         "line1\nline2\nline3",
			wantLines:    3,
			wantFrags:    1,
		},
		{
			name:         "default fragment size when zero",
			fragmentSize: 0,
			text:         "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\nline11",
			wantLines:    11,
			wantFrags:    2, // 10 + 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewTextProcessor(tt.fragmentSize)
			lines, frags := p.ProcessText(tt.text)
			if lines != tt.wantLines {
				t.Errorf("ProcessText() lines = %v, want %v", lines, tt.wantLines)
			}
			if len(frags) != tt.wantFrags {
				t.Errorf("ProcessText() fragments = %v, want %v", len(frags), tt.wantFrags)
			}
			// Verify total lines match sum of fragment sizes
			totalFragLines := 0
			for _, frag := range frags {
				totalFragLines += len(frag)
			}
			if totalFragLines != tt.wantLines {
				t.Errorf("ProcessText() total fragment lines = %v, want %v", totalFragLines, tt.wantLines)
			}
		})
	}
}
