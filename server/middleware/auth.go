package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"opm/models"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey contextKey = "user"

// RequireAuthMiddleware ensures the request has a valid JWT token
func RequireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get JWT secret from environment
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}
		// Get token from cookie or Authorization header
		token := ""

		// Check cookie first
		if cookie, err := r.Cookie("token"); err == nil {
			token = cookie.Value
		}

		// Check Authorization header as fallback
		if token == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse and validate token
		claims := &Claims{}
		jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !jwtToken.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		authUser := &models.AuthUser{
			UserID: claims.UserID,
			Token:  token,
		}
		ctx := context.WithValue(r.Context(), userContextKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware extracts auth info if present but doesn't require it
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get JWT secret from environment
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			// Continue without auth if not configured
			next.ServeHTTP(w, r)
			return
		}
		
		// Get token from cookie or Authorization header
		token := ""

		// Check cookie first
		if cookie, err := r.Cookie("token"); err == nil {
			token = cookie.Value
		}

		// Check Authorization header as fallback
		if token == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// If no token, continue without auth
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Parse and validate token
		claims := &Claims{}
		jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		// If token is invalid, continue without auth
		if err != nil || !jwtToken.Valid {
			next.ServeHTTP(w, r)
			return
		}

		// Add user info to context
		authUser := &models.AuthUser{
			UserID: claims.UserID,
			Token:  token,
		}
		ctx := context.WithValue(r.Context(), userContextKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetAuthUser retrieves the authenticated user from context
func GetAuthUser(ctx context.Context) (*models.AuthUser, bool) {
	user, ok := ctx.Value(userContextKey).(*models.AuthUser)
	return user, ok
}

// Claims represents JWT claims
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
