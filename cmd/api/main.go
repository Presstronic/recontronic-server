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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/presstronic/recontronic-server/internal/config"
	"github.com/presstronic/recontronic-server/internal/database"
	"github.com/presstronic/recontronic-server/internal/handlers"
	authmiddleware "github.com/presstronic/recontronic-server/internal/middleware"
	"github.com/presstronic/recontronic-server/internal/repository"
	"github.com/presstronic/recontronic-server/internal/services"
	"github.com/presstronic/recontronic-server/pkg/validator"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Recontronic API Server...")
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("REST API on port: %d", cfg.Server.RESTPort)

	// Initialize database connection
	db, err := database.NewPostgresConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)
	log.Printf("✓ Connected to database")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)

	// Initialize validator
	v := validator.New()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, v)

	// Initialize middleware
	authMiddleware := authmiddleware.NewAuthMiddleware(authService)

	// Setup router
	r := setupRouter(cfg, authHandler, authMiddleware, db)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.RESTPort),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("✓ REST API server listening on :%d", cfg.Server.RESTPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✓ Server stopped gracefully")
}

func setupRouter(cfg *config.Config, authHandler *handlers.AuthHandler, authMiddleware *authmiddleware.AuthMiddleware, db interface{ Ping() error }) chi.Router {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Security.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoints (no auth required)
	r.Get("/health", healthCheckHandler)
	r.Get("/ready", readinessCheckHandler(db))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public authentication routes
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			// Auth routes
			r.Get("/auth/me", authHandler.Me)
			r.Post("/auth/keys", authHandler.CreateAPIKey)
			r.Get("/auth/keys", authHandler.ListAPIKeys)
			r.Delete("/auth/keys/{id}", authHandler.RevokeAPIKey)

			// TODO: Add program routes
			// TODO: Add scan routes
			// TODO: Add anomaly routes
		})
	})

	return r
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func readinessCheckHandler(db interface{ Ping() error }) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"not ready","error":"database connection failed"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	}
}
