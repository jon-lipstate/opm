package main

import (
	"context"
	"log"
	"net/http"
	"opm/config"
	"opm/db"
	"opm/handlers/auth"
	"opm/handlers/packages"
	"opm/handlers/tags"
	"opm/handlers/users"
	"opm/middleware"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	// Initialize database connection pool
	db.InitPool(cfg.DatabaseURL)
	defer db.Close()

	// No migrations - database setup is done manually

	// Create router
	r := mux.NewRouter()

	// API routes
	api := r

	// Apply global middleware
	api.Use(middleware.Logger)
	api.Use(middleware.RequestID)
	api.Use(middleware.RateLimit(cfg.RateLimit, cfg.RateWindow))

	// Authenticated routes subrouter
	authApi := api.NewRoute().Subrouter()
	authApi.Use(middleware.RequireAuthMiddleware)

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Auth routes (these don't require authentication)
	api.HandleFunc("/auth/github", auth.GitHubLogin(cfg)).Methods("GET")
	api.HandleFunc("/auth/github/callback", auth.GitHubCallback(cfg)).Methods("GET")
	api.HandleFunc("/auth/discord", auth.DiscordLogin(cfg)).Methods("GET")
	api.HandleFunc("/auth/discord/callback", auth.DiscordCallback(cfg)).Methods("GET")
	api.HandleFunc("/auth/logout", auth.Logout()).Methods("POST")
	authApi.HandleFunc("/auth/me", users.GetCurrentUser).Methods("GET")

	// Public routes with optional auth (for bookmark/vote status)
	optionalAuthApi := api.NewRoute().Subrouter()
	optionalAuthApi.Use(middleware.OptionalAuthMiddleware)

	// Package routes
	optionalAuthApi.HandleFunc("/packages", packages.List).Methods("GET")
	optionalAuthApi.HandleFunc("/packages/search", packages.Search).Methods("GET")
	authApi.HandleFunc("/packages", packages.Create).Methods("POST")
	authApi.HandleFunc("/packages/{id}", packages.Update).Methods("PUT")
	authApi.HandleFunc("/packages/{id}", packages.Delete).Methods("DELETE")
	optionalAuthApi.HandleFunc("/packages/{alias}/{slug}", packages.Get).Methods("GET")
	optionalAuthApi.HandleFunc("/readme", packages.GetPackageReadme).Methods("GET")
	authApi.HandleFunc("/repository/metadata", packages.GetRepositoryMetadata).Methods("GET")

	// Bookmark routes (require auth)
	authApi.HandleFunc("/packages/{alias}/{slug}/bookmark", packages.Bookmark).Methods("POST")
	authApi.HandleFunc("/packages/{alias}/{slug}/bookmark", packages.Unbookmark).Methods("DELETE")

	// Tag routes (require auth)
	authApi.HandleFunc("/tags", packages.AddTag).Methods("POST")
	authApi.HandleFunc("/tags/vote", packages.VoteTag).Methods("POST")

	// Flag/moderation routes
	authApi.HandleFunc("/flags", packages.FlagPackage).Methods("POST")
	optionalAuthApi.HandleFunc("/flags", packages.GetPackageFlags).Methods("GET")
	optionalAuthApi.HandleFunc("/flags/stats", packages.GetFlagStats).Methods("GET")
	authApi.HandleFunc("/flags/all", packages.GetAllFlags).Methods("GET")          // Moderator only
	authApi.HandleFunc("/flags/{id}/resolve", packages.ResolveFlag).Methods("PUT") // Moderator only
	authApi.HandleFunc("/flags/{id}", packages.DeleteFlag).Methods("DELETE")
	authApi.HandleFunc("/users/me/flags", packages.GetUserFlags).Methods("GET")

	// Tags
	api.HandleFunc("/tags", tags.List).Methods("GET")

	// User routes
	authApi.HandleFunc("/users/me/packages", users.ListUserPackages).Methods("GET")
	authApi.HandleFunc("/users/me", users.UpdateProfile).Methods("PUT")
	authApi.HandleFunc("/users/check-alias", users.CheckAliasAvailability).Methods("GET")
	// authApi.HandleFunc("/users/me/bookmarks", users.ListBookmarks).Methods("GET")

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:9000", "https://pkg-odin.org/"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Cookie"},
		ExposedHeaders:   []string{"Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(r)

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
