package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type requestIDKey string

const requestIDContextKey requestIDKey = "request_id"

// RequestID middleware adds a unique request ID to each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate request ID
		id := generateRequestID()
		
		// Add to context
		ctx := context.WithValue(r.Context(), requestIDContextKey, id)
		
		// Add to response header
		w.Header().Set("X-Request-ID", id)
		
		// Continue
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDContextKey).(string); ok {
		return id
	}
	return ""
}

func generateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}