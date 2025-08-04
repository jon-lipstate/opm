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
	"opm/logger"
	"opm/middleware"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// load .env file
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.InitLoggers()
	mainLogger := logger.MainLogger

	mainLogger.Printf("🚀 Starting OPM server in %s mode", cfg.Env)

	db.InitPool(cfg.DatabaseURL)
	defer db.Close()

	r := mux.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RateLimit(cfg.RateLimit, cfg.RateWindow))

	authApi := r.NewRoute().Subrouter()
	authApi.Use(middleware.RequireAuthMiddleware)

	// Public routes with optional auth (for bookmark/vote status)
	optionalAuthApi := r.NewRoute().Subrouter()
	optionalAuthApi.Use(middleware.OptionalAuthMiddleware)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Auth routes (these don't require authentication)
	r.HandleFunc("/auth/github", auth.GitHubLogin(cfg)).Methods("GET")
	r.HandleFunc("/auth/github/callback", auth.GitHubCallback(cfg)).Methods("GET")
	r.HandleFunc("/auth/discord", auth.DiscordLogin(cfg)).Methods("GET")
	r.HandleFunc("/auth/discord/callback", auth.DiscordCallback(cfg)).Methods("GET")
	r.HandleFunc("/auth/logout", auth.Logout()).Methods("POST")
	authApi.HandleFunc("/auth/me", users.GetCurrentUser).Methods("GET")

	// Package routes
	optionalAuthApi.HandleFunc("/readme", packages.GetPackageReadme).Methods("GET")
	authApi.HandleFunc("/repository/metadata", packages.GetRepositoryMetadata).Methods("GET")
	optionalAuthApi.HandleFunc("/packages", packages.List).Methods("GET")
	optionalAuthApi.HandleFunc("/packages/search", packages.Search).Methods("GET")
	authApi.HandleFunc("/packages", packages.Create).Methods("POST")
	authApi.HandleFunc("/packages/bookmark", packages.Bookmark).Methods("POST")     // param: package_id
	authApi.HandleFunc("/packages/bookmark", packages.Unbookmark).Methods("DELETE") // param: package_id
	// MUST BE BELOW OTHER ROUTES DUE TO WILDCARD MUX:
	optionalAuthApi.HandleFunc("/packages/{userSlug}/{pkgSlug}", packages.Get).Methods("GET")
	authApi.HandleFunc("/packages/{id}", packages.Update).Methods("PUT")
	authApi.HandleFunc("/packages/{id}", packages.Delete).Methods("DELETE")

	// Tag routes (require auth)
	authApi.HandleFunc("/tags", packages.AddTag).Methods("POST")
	authApi.HandleFunc("/tags/vote", packages.VoteTag).Methods("POST") // param: package_id

	// Flag/moderation routes
	authApi.HandleFunc("/flags", packages.FlagPackage).Methods("POST")
	optionalAuthApi.HandleFunc("/flags", packages.GetPackageFlags).Methods("GET")    // param: package_id
	optionalAuthApi.HandleFunc("/flags/stats", packages.GetFlagStats).Methods("GET") // param: package_id
	authApi.HandleFunc("/flags/all", packages.GetAllFlags).Methods("GET")            // Moderator only
	authApi.HandleFunc("/users/me/flags", packages.GetUserFlags).Methods("GET")
	authApi.HandleFunc("/flags/{id}/resolve", packages.ResolveFlag).Methods("PUT") // Moderator only
	authApi.HandleFunc("/flags/{id}", packages.DeleteFlag).Methods("DELETE")

	// Tags
	r.HandleFunc("/tags", tags.List).Methods("GET")

	// User routes
	authApi.HandleFunc("/users/me/packages", users.ListUserPackages).Methods("GET")
	authApi.HandleFunc("/users/me", users.UpdateProfile).Methods("PUT")
	authApi.HandleFunc("/users/check-user-slug", users.CheckSlugAvailability).Methods("GET")
	// authApi.HandleFunc("/users/me/bookmarks", users.ListBookmarks).Methods("GET")

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:9000", "https://pkg-odin.org"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Cookie", "Content-Disposition"},
		ExposedHeaders:   []string{"Set-Cookie", "Content-Disposition"},
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
		mainLogger.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			mainLogger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	mainLogger.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		mainLogger.Fatalf("Server forced to shutdown: %v", err)
	}

	mainLogger.Println("Server exited")
}
