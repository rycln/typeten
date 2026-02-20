package handlers

import "typeten/internal/domain"

// CreateTextRequest represents the HTTP request for creating a text.
type CreateTextRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreateTextResponse represents the HTTP response for creating a text.
type CreateTextResponse struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	TotalLines    int    `json:"total_lines"`
	FragmentSize  int    `json:"fragment_size"`
	FragmentCount int    `json:"fragment_count"`
	CreatedAt     string `json:"created_at"`
}

// CreateSessionRequest represents the HTTP request for creating a session.
type CreateSessionRequest struct {
	TextID string `json:"text_id"`
}

// CreateSessionResponse represents the HTTP response for creating a session.
type CreateSessionResponse struct {
	ID                   string  `json:"id"`
	UserID               string  `json:"user_id"`
	TextID               string  `json:"text_id"`
	CurrentFragmentIdx   int     `json:"current_fragment_idx"`
	CurrentLineIdx       int     `json:"current_line_idx"`
	CompletedLines       int     `json:"completed_lines"`
	TotalAccuracyPercent float64 `json:"total_accuracy_percent"`
	AverageWPM           float64 `json:"average_wpm"`
	IsCompleted          bool    `json:"is_completed"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

// RecordProgressRequest represents the HTTP request for recording progress.
type RecordProgressRequest struct {
	AccuracyPercent float64 `json:"accuracy_percent"`
	WPM             float64 `json:"wpm"`
}

// RecordProgressResponse represents the HTTP response for recording progress.
type RecordProgressResponse struct {
	ID                   string  `json:"id"`
	CompletedLines       int     `json:"completed_lines"`
	TotalAccuracyPercent float64 `json:"total_accuracy_percent"`
	AverageWPM           float64 `json:"average_wpm"`
	IsCompleted          bool    `json:"is_completed"`
}

// GetSessionResponse represents the HTTP response for getting a session.
type GetSessionResponse struct {
	ID                   string  `json:"id"`
	UserID               string  `json:"user_id"`
	TextID               string  `json:"text_id"`
	CurrentFragmentIdx   int     `json:"current_fragment_idx"`
	CurrentLineIdx       int     `json:"current_line_idx"`
	CompletedLines       int     `json:"completed_lines"`
	TotalAccuracyPercent float64 `json:"total_accuracy_percent"`
	AverageWPM           float64 `json:"average_wpm"`
	IsCompleted          bool    `json:"is_completed"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

// ListTextsResponse represents the HTTP response for listing texts.
type ListTextsResponse struct {
	Texts []TextInfoResponse `json:"texts"`
}

// TextInfoResponse represents a text info in responses.
type TextInfoResponse struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	TotalLines    int    `json:"total_lines"`
	FragmentSize  int    `json:"fragment_size"`
	FragmentCount int    `json:"fragment_count"`
	CreatedAt     string `json:"created_at"`
}

// GetTextFragmentsResponse represents the HTTP response for getting fragments.
type GetTextFragmentsResponse struct {
	Fragments []FragmentResponse `json:"fragments"`
}

// FragmentResponse represents a fragment in responses.
type FragmentResponse struct {
	ID          string   `json:"id"`
	FragmentIdx int     `json:"fragment_idx"`
	Lines       []string `json:"lines"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Helper functions to convert domain models to DTOs
func textInfoToResponse(info *domain.TextInfo) TextInfoResponse {
	return TextInfoResponse{
		ID:            string(info.ID),
		Title:         info.Title,
		TotalLines:    info.TotalLines,
		FragmentSize:  info.FragmentSize,
		FragmentCount: info.FragmentCount,
		CreatedAt:     info.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func sessionToResponse(session *domain.Session) GetSessionResponse {
	return GetSessionResponse{
		ID:                   string(session.ID),
		UserID:               string(session.UserID),
		TextID:               string(session.TextID),
		CurrentFragmentIdx:   session.CurrentFragmentIdx,
		CurrentLineIdx:       session.CurrentLineIdx,
		CompletedLines:       session.CompletedLines,
		TotalAccuracyPercent: session.TotalAccuracyPercent,
		AverageWPM:           session.AverageWPM,
		IsCompleted:          session.IsCompleted,
		CreatedAt:            session.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:            session.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
