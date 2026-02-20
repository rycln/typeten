package handlers

import (
	"net/http"
	"strings"
)

// Router sets up HTTP routes for the API.
type Router struct {
	handlers *Handlers
}

// NewRouter creates a new router.
func NewRouter(handlers *Handlers) *Router {
	return &Router{handlers: handlers}
}

// ServeHTTP implements http.Handler and routes requests to appropriate handlers.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/" && r.Method == http.MethodGet:
		rt.handlers.IndexPage(w, r)
	case path == "/texts" && r.Method == http.MethodPost:
		rt.handlers.CreateTextHTML(w, r)
	case strings.HasPrefix(path, "/texts/") && r.Method == http.MethodGet:
		// /texts/{id}
		id := strings.TrimPrefix(path, "/texts/")
		if id == "" {
			http.NotFound(w, r)
			return
		}
		rt.handlers.TextDetailPage(w, r, id)
	case path == "/sessions" && r.Method == http.MethodPost:
		rt.handlers.CreateSessionHTML(w, r)
	case strings.HasPrefix(path, "/sessions/") && r.Method == http.MethodGet:
		// /sessions/{id}
		id := strings.TrimPrefix(path, "/sessions/")
		if id == "" {
			http.NotFound(w, r)
			return
		}
		rt.handlers.SessionPage(w, r, id)
	case path == "/api/texts" && r.Method == http.MethodPost:
		rt.handlers.CreateText(w, r)
	case path == "/api/texts" && r.Method == http.MethodGet:
		rt.handlers.ListTexts(w, r)
	case strings.HasPrefix(path, "/api/texts/") && strings.HasSuffix(path, "/fragments") && r.Method == http.MethodGet:
		rt.handlers.GetTextFragments(w, r)
	case path == "/api/sessions" && r.Method == http.MethodPost:
		rt.handlers.CreateSession(w, r)
	case strings.HasPrefix(path, "/api/sessions/") && strings.HasSuffix(path, "/progress") && r.Method == http.MethodPost:
		rt.handlers.RecordProgress(w, r)
	case strings.HasPrefix(path, "/api/sessions/") && !strings.HasSuffix(path, "/progress") && r.Method == http.MethodGet:
		rt.handlers.GetSession(w, r)
	default:
		http.NotFound(w, r)
	}
}
