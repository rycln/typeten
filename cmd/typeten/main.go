package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"typeten/internal/domain"
	"typeten/internal/handlers"
	infraRepo "typeten/internal/infrastructure/repository"
	"typeten/internal/repository"
	"typeten/internal/usecases"
)

const (
	defaultPort        = "8080"
	defaultFragmentSize = 10
)

func main() {
	// Initialize repositories (in-memory for MVP)
	userRepo := infraRepo.NewMemoryUserRepository()
	textRepo := infraRepo.NewMemoryTextRepository()
	sessionRepo := infraRepo.NewMemorySessionRepository()

	// Create a default user for MVP (in production, this would come from auth)
	ctx := context.Background()
	defaultUser, err := createDefaultUser(ctx, userRepo)
	if err != nil {
		log.Fatalf("Failed to create default user: %v", err)
	}

	// Initialize use cases
	createTextUseCase := usecases.NewCreateTextUseCase(textRepo, userRepo, defaultFragmentSize)
	createSessionUseCase := usecases.NewCreateSessionUseCase(sessionRepo, textRepo, userRepo)
	recordProgressUseCase := usecases.NewRecordProgressUseCase(sessionRepo)
	getSessionUseCase := usecases.NewGetSessionUseCase(sessionRepo)
	listTextsUseCase := usecases.NewListTextsUseCase(textRepo, userRepo)
	getTextFragmentsUseCase := usecases.NewGetTextFragmentsUseCase(textRepo)

	// Initialize handlers
	httpHandlers := handlers.NewHandlers(
		createTextUseCase,
		createSessionUseCase,
		recordProgressUseCase,
		getSessionUseCase,
		listTextsUseCase,
		getTextFragmentsUseCase,
		defaultUser.ID,
	)

	// Setup router
	router := handlers.NewRouter(httpHandlers)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		log.Printf("API endpoints:")
		log.Printf("  POST   /api/texts")
		log.Printf("  GET    /api/texts")
		log.Printf("  GET    /api/texts/:id/fragments")
		log.Printf("  POST   /api/sessions")
		log.Printf("  GET    /api/sessions/:id")
		log.Printf("  POST   /api/sessions/:id/progress")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// createDefaultUser creates a default user for MVP testing.
func createDefaultUser(ctx context.Context, userRepo repository.UserRepository) (*domain.User, error) {
	// Try to get existing user
	userID := domain.UserID("user_default")
	user, err := userRepo.GetByID(ctx, userID)
	if err == nil {
		return user, nil
	}

	// Create new user
	now := time.Now()
	user, err = domain.NewUser(
		userID,
		"user@example.com",
		"default_user",
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to store user: %w", err)
	}

	return user, nil
}
