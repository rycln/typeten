package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"typeten/internal/domain"
	"typeten/internal/usecases"
)

// Handlers holds all HTTP handlers and their dependencies.
type Handlers struct {
	createTextUseCase        *usecases.CreateTextUseCase
	createSessionUseCase     *usecases.CreateSessionUseCase
	recordProgressUseCase    *usecases.RecordProgressUseCase
	getSessionUseCase        *usecases.GetSessionUseCase
	listTextsUseCase         *usecases.ListTextsUseCase
	getTextFragmentsUseCase  *usecases.GetTextFragmentsUseCase
	currentUserID            domain.UserID // MVP: single user, will be replaced with auth
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(
	createTextUseCase *usecases.CreateTextUseCase,
	createSessionUseCase *usecases.CreateSessionUseCase,
	recordProgressUseCase *usecases.RecordProgressUseCase,
	getSessionUseCase *usecases.GetSessionUseCase,
	listTextsUseCase *usecases.ListTextsUseCase,
	getTextFragmentsUseCase *usecases.GetTextFragmentsUseCase,
	currentUserID domain.UserID,
) *Handlers {
	return &Handlers{
		createTextUseCase:       createTextUseCase,
		createSessionUseCase:    createSessionUseCase,
		recordProgressUseCase:   recordProgressUseCase,
		getSessionUseCase:       getSessionUseCase,
		listTextsUseCase:        listTextsUseCase,
		getTextFragmentsUseCase: getTextFragmentsUseCase,
		currentUserID:           currentUserID,
	}
}

// CreateText handles POST /api/texts
func (h *Handlers) CreateText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	input := usecases.CreateTextInput{
		UserID:  h.currentUserID,
		Title:   req.Title,
		Content: req.Content,
	}

	output, err := h.createTextUseCase.Execute(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create text: %v", err))
		return
	}

	resp := CreateTextResponse{
		ID:            string(output.TextInfo.ID),
		Title:         output.TextInfo.Title,
		TotalLines:    output.TextInfo.TotalLines,
		FragmentSize:  output.TextInfo.FragmentSize,
		FragmentCount: output.TextInfo.FragmentCount,
		CreatedAt:     output.TextInfo.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	respondJSON(w, http.StatusCreated, resp)
}

// CreateSession handles POST /api/sessions
func (h *Handlers) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	input := usecases.CreateSessionInput{
		UserID: h.currentUserID,
		TextID: domain.TextID(req.TextID),
	}

	output, err := h.createSessionUseCase.Execute(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create session: %v", err))
		return
	}

	resp := sessionToResponse(output.Session)
	respondJSON(w, http.StatusCreated, resp)
}

// RecordProgress handles POST /api/sessions/:id/progress
func (h *Handlers) RecordProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract session ID from path like "/api/sessions/abc123/progress"
	path := r.URL.Path
	sessionID := strings.TrimPrefix(path, "/api/sessions/")
	sessionID = strings.TrimSuffix(sessionID, "/progress")

	var req RecordProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	input := usecases.RecordProgressInput{
		SessionID:       sessionID,
		AccuracyPercent: req.AccuracyPercent,
		WPM:             req.WPM,
	}

	output, err := h.recordProgressUseCase.Execute(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to record progress: %v", err))
		return
	}

	resp := RecordProgressResponse{
		ID:                   string(output.Session.ID),
		CompletedLines:       output.Session.CompletedLines,
		TotalAccuracyPercent: output.Session.TotalAccuracyPercent,
		AverageWPM:           output.Session.AverageWPM,
		IsCompleted:          output.Session.IsCompleted,
	}

	respondJSON(w, http.StatusOK, resp)
}

// GetSession handles GET /api/sessions/:id
func (h *Handlers) GetSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Path[len("/api/sessions/"):]

	input := usecases.GetSessionInput{
		SessionID: domain.SessionID(sessionID),
	}

	output, err := h.getSessionUseCase.Execute(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Session not found: %v", err))
		return
	}

	resp := sessionToResponse(output.Session)
	respondJSON(w, http.StatusOK, resp)
}

// ListTexts handles GET /api/texts
func (h *Handlers) ListTexts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	input := usecases.ListTextsInput{
		UserID: h.currentUserID,
	}

	output, err := h.listTextsUseCase.Execute(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list texts: %v", err))
		return
	}

	texts := make([]TextInfoResponse, len(output.Texts))
	for i, text := range output.Texts {
		texts[i] = textInfoToResponse(text)
	}

	resp := ListTextsResponse{Texts: texts}
	respondJSON(w, http.StatusOK, resp)
}

// GetTextFragments handles GET /api/texts/:id/fragments
func (h *Handlers) GetTextFragments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract text ID from path like "/api/texts/abc123/fragments"
	path := r.URL.Path
	textID := strings.TrimPrefix(path, "/api/texts/")
	textID = strings.TrimSuffix(textID, "/fragments")

	input := usecases.GetTextFragmentsInput{
		TextID: domain.TextID(textID),
	}

	output, err := h.getTextFragmentsUseCase.Execute(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Failed to get fragments: %v", err))
		return
	}

	fragments := make([]FragmentResponse, len(output.Fragments))
	for i, frag := range output.Fragments {
		fragments[i] = FragmentResponse{
			ID:          string(frag.ID),
			FragmentIdx: frag.FragmentIdx,
			Lines:       frag.Lines(),
		}
	}

	resp := GetTextFragmentsResponse{Fragments: fragments}
	respondJSON(w, http.StatusOK, resp)
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
