package usecases

import "strings"

// TextProcessor splits raw text content into fragments of a specified size.
type TextProcessor struct {
	FragmentSize int
}

// NewTextProcessor creates a new text processor with the given fragment size.
func NewTextProcessor(fragmentSize int) *TextProcessor {
	if fragmentSize <= 0 {
		fragmentSize = 10 // default
	}
	return &TextProcessor{FragmentSize: fragmentSize}
}

// ProcessText splits text into lines and then into fragments.
// Returns the total line count, fragment count, and the fragments themselves.
func (p *TextProcessor) ProcessText(text string) (totalLines int, fragments [][]string) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return 0, nil
	}
	lines := strings.Split(trimmed, "\n")
	totalLines = len(lines)
	if totalLines == 0 {
		return 0, nil
	}
	
	// Split lines into fragments
	for i := 0; i < totalLines; i += p.FragmentSize {
		end := i + p.FragmentSize
		if end > totalLines {
			end = totalLines
		}
		fragments = append(fragments, lines[i:end])
	}
	
	return totalLines, fragments
}
