package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"typeten/internal/domain"
	"typeten/internal/usecases"
)

func setupTestHandlers(t *testing.T) *Handlers {
	userRepo := usecases.NewMockUserRepository()
	textRepo := usecases.NewMockTextRepository()
	sessionRepo := usecases.NewMockSessionRepository()

	now := time.Now()
	user, err := domain.NewUser("user_1", "test@example.com", "testuser", now)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := userRepo.Create(context.Background(), user); err != nil {
		t.Fatalf("Failed to store user: %v", err)
	}

	createTextUseCase := usecases.NewCreateTextUseCase(textRepo, userRepo, 5)
	createSessionUseCase := usecases.NewCreateSessionUseCase(sessionRepo, textRepo, userRepo)
	recordProgressUseCase := usecases.NewRecordProgressUseCase(sessionRepo)
	getSessionUseCase := usecases.NewGetSessionUseCase(sessionRepo)
	listTextsUseCase := usecases.NewListTextsUseCase(textRepo, userRepo)
	getTextFragmentsUseCase := usecases.NewGetTextFragmentsUseCase(textRepo)

	return NewHandlers(
		createTextUseCase,
		createSessionUseCase,
		recordProgressUseCase,
		getSessionUseCase,
		listTextsUseCase,
		getTextFragmentsUseCase,
		user.ID,
	)
}

func TestHandlers_CreateText(t *testing.T) {
	handlers := setupTestHandlers(t)

	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{
		{
			name:       "valid request",
			method:     http.MethodPost,
			body:       `{"title":"Test Text","content":"line1\nline2\nline3"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid JSON",
			method:     http.MethodPost,
			body:       `invalid json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "wrong method",
			method:     http.MethodGet,
			body:       `{"title":"Test","content":"line1"}`,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/texts", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handlers.CreateText(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateText() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandlers_ListTexts(t *testing.T) {
	handlers := setupTestHandlers(t)

	// Create a text first using the handlers' use case
	ctx := context.Background()
	_, err := handlers.createTextUseCase.Execute(ctx, usecases.CreateTextInput{
		UserID:  handlers.currentUserID,
		Title:   "Test Text",
		Content: "line1\nline2",
	})
	if err != nil {
		t.Fatalf("Failed to create test text: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/texts", nil)
	w := httptest.NewRecorder()

	handlers.ListTexts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("ListTexts() status = %v, want %v", w.Code, http.StatusOK)
	}

	var resp ListTextsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(resp.Texts) == 0 {
		t.Error("ListTexts() returned empty list")
	}
}

func TestHandlers_CreateSession(t *testing.T) {
	handlers := setupTestHandlers(t)

	// Create a text first using the handlers' use case
	ctx := context.Background()
	output, err := handlers.createTextUseCase.Execute(ctx, usecases.CreateTextInput{
		UserID:  handlers.currentUserID,
		Title:   "Test Text",
		Content: "line1\nline2",
	})
	if err != nil {
		t.Fatalf("Failed to create test text: %v", err)
	}

	body := `{"text_id":"` + string(output.TextInfo.ID) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/sessions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlers.CreateSession(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("CreateSession() status = %v, want %v", w.Code, http.StatusCreated)
	}

	var resp CreateSessionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.TextID != string(output.TextInfo.ID) {
		t.Errorf("CreateSession() TextID = %v, want %v", resp.TextID, output.TextInfo.ID)
	}
}

func TestHandlers_RecordProgress(t *testing.T) {
	handlers := setupTestHandlers(t)

	// Create a text and session first
	ctx := context.Background()
	textOutput, err := handlers.createTextUseCase.Execute(ctx, usecases.CreateTextInput{
		UserID:  handlers.currentUserID,
		Title:   "Test Text",
		Content: "line1\nline2",
	})
	if err != nil {
		t.Fatalf("Failed to create test text: %v", err)
	}

	sessionOutput, err := handlers.createSessionUseCase.Execute(ctx, usecases.CreateSessionInput{
		UserID: handlers.currentUserID,
		TextID: textOutput.TextInfo.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}

	body := `{"accuracy_percent":95.5,"wpm":45.2}`
	req := httptest.NewRequest(http.MethodPost, "/api/sessions/"+string(sessionOutput.Session.ID)+"/progress", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlers.RecordProgress(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("RecordProgress() status = %v, want %v", w.Code, http.StatusOK)
	}

	var resp RecordProgressResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.CompletedLines != 1 {
		t.Errorf("RecordProgress() CompletedLines = %v, want 1", resp.CompletedLines)
	}
}

func TestHandlers_GetSession(t *testing.T) {
	handlers := setupTestHandlers(t)

	// Create a text and session first
	ctx := context.Background()
	textOutput, err := handlers.createTextUseCase.Execute(ctx, usecases.CreateTextInput{
		UserID:  handlers.currentUserID,
		Title:   "Test Text",
		Content: "line1\nline2",
	})
	if err != nil {
		t.Fatalf("Failed to create test text: %v", err)
	}

	sessionOutput, err := handlers.createSessionUseCase.Execute(ctx, usecases.CreateSessionInput{
		UserID: handlers.currentUserID,
		TextID: textOutput.TextInfo.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/sessions/"+string(sessionOutput.Session.ID), nil)
	w := httptest.NewRecorder()

	handlers.GetSession(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetSession() status = %v, want %v", w.Code, http.StatusOK)
	}

	var resp GetSessionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.ID != string(sessionOutput.Session.ID) {
		t.Errorf("GetSession() ID = %v, want %v", resp.ID, sessionOutput.Session.ID)
	}
}

func TestHandlers_GetTextFragments(t *testing.T) {
	handlers := setupTestHandlers(t)

	// Create a text with fragments first using the handlers' use case
	ctx := context.Background()
	output, err := handlers.createTextUseCase.Execute(ctx, usecases.CreateTextInput{
		UserID:  handlers.currentUserID,
		Title:   "Test Text",
		Content: "line1\nline2\nline3",
	})
	if err != nil {
		t.Fatalf("Failed to create test text: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/texts/"+string(output.TextInfo.ID)+"/fragments", nil)
	w := httptest.NewRecorder()

	handlers.GetTextFragments(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetTextFragments() status = %v, want %v", w.Code, http.StatusOK)
	}

	var resp GetTextFragmentsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(resp.Fragments) == 0 {
		t.Error("GetTextFragments() returned empty fragments")
	}
}

func TestRouter(t *testing.T) {
	handlers := setupTestHandlers(t)
	router := NewRouter(handlers)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "index page",
			method:     http.MethodGet,
			path:       "/",
			wantStatus: http.StatusOK,
		},
		{
			name:       "create text form - needs setup",
			method:     http.MethodPost,
			path:       "/texts",
			wantStatus: http.StatusSeeOther, // Will redirect on success
		},
		{
			name:       "create text form",
			method:     http.MethodPost,
			path:       "/texts",
			wantStatus: http.StatusSeeOther,
		},
		{
			name:       "api create text",
			method:     http.MethodPost,
			path:       "/api/texts",
			wantStatus: http.StatusBadRequest, // Missing body
		},
		{
			name:       "api list texts",
			method:     http.MethodGet,
			path:       "/api/texts",
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			method:     http.MethodGet,
			path:       "/nonexistent",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.method == http.MethodPost && tt.path == "/texts" {
				body = bytes.NewBufferString("title=Test&content=line1")
			}
			req := httptest.NewRequest(tt.method, tt.path, body)
			if tt.method == http.MethodPost && tt.path == "/texts" {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Router.ServeHTTP() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
